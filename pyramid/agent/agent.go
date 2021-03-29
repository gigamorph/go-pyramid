package agent

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gigamorph/go-pyramid/config"
	"github.com/gigamorph/go-pyramid/pyramid/context"
	"github.com/gigamorph/go-pyramid/pyramid/input"
	"github.com/gigamorph/go-pyramid/pyramid/output"
	im "github.com/gigamorph/go-pyramid/shellcmds/imagemagick"
	"github.com/gigamorph/go-pyramid/shellcmds/tiff"
	"github.com/gigamorph/go-pyramid/shellcmds/vips"
	"github.com/gigamorph/go-pyramid/util"
)

// Agent is a wrapper around tools and operations used to convert
// images to pyramidal TIFF.
//
// There should be only one instance of Agent running at any time.
// A typical usage is:
//   agent := NewAgent()
//   agent.Initialize()
//   defer agent.Finalize()
//   agent.Convert(params) // params is of convert.Params type
type Agent struct {
}

// New returns a new instance of Agent.
func New() *Agent {
	agent := Agent{}
	return &agent
}

// Convert is the public method to call to actually convert an image.
// p contains input, output, and other information needed for conversion.
func (a *Agent) Convert(p input.Params) (*output.Params, error) {
	c := context.New(p)
	a.mkdirp(c.Input.TempDir)

	err := a.toPyramidTIFF(c)
	if err != nil {
		return nil, fmt.Errorf("pyramid.agent.Agent#Convert failed to create pyramid - %v", err)
	}
	if p.DeleteTemp {
		err = os.RemoveAll(c.Input.TempDir)
		if err != nil {
			log.Printf("ERROR pyramid.agent.Agent#Convert failed to delete temp dir %s - %v\n", c.Input.TempDir, err)
		}
	}
	return &c.Output, nil
}

func (a *Agent) toPyramidTIFF(c *context.Context) (err error) {
	targetICCProfile := config.TargetICCProfileIIIF
	if c.Input.TargetICCProfile != "" {
		targetICCProfile = c.Input.TargetICCProfile
	}

	// Make sure input is a single file TIFF
	if err = vips.ToTiff(fmt.Sprintf("%s[0]", c.Input.InFile), c.TiffFile); err != nil {
		return fmt.Errorf("pyramid.agent.Agent#ToPyramidTIFF failed to convert %s to TIFF - %v", c.Input.InFile, err)
	}
	if err != nil {
		return fmt.Errorf("pyramid.agent.Agent#ToPyramidTIFF failed to convert to tiff - %v", err)
	}

	tiff := c.TiffFile

	c.Width, err = vips.Width(c.TiffFile)
	if err != nil {
		return fmt.Errorf("pyramid.agent.Agent#ToPyramidTIFF failed to get width - %v", err)
	}
	c.Height, err = vips.Height(c.TiffFile)
	if err != nil {
		return fmt.Errorf("pyramid.agent.Agent#ToPyramidTIFF failed to get height - %v", err)
	}

	c.Output.InputWidth = c.Width
	c.Output.InputHeight = c.Height

	imageFormat, channels, depth, iccProfileName, err := im.GetInfo(tiff)
	if err != nil {
		return fmt.Errorf("pyramid.agent.Agent#toPyramidTIFF failed get info from %s - %v", tiff, err)
	}
	depth64, err := strconv.ParseUint(depth, 10, 64)
	if err != nil {
		return fmt.Errorf("pyramid.agent.Agent#toPyramidTIFF failed to parse depth - %v", err)
	}
	c.BitDepth = uint(depth64)

	log.Printf("imageFormat: %s, channels: %s, profile: %s\n", imageFormat, channels, iccProfileName)

	// Check if channels is supported
	if valid := a.validateChannels(channels); !valid {
		return fmt.Errorf("Image %s has channels %s which is not supported at this time",
			tiff, channels)
	}

	// We have to flatten the image to remove the alpha channel / trasparency
	// before proceeding
	if channels == "srgba" {
		if err = vips.RemoveAlpha(tiff, c.NoalphaFile); err != nil {
			return fmt.Errorf("Agent#toPyramidTIFF RemoveAlpha failed - %v", err)
		}
	} else {
		c.NoalphaFile = tiff
	}

	newProfile := false

	// In the case of gray with no embedded color profile or with an embedded
	// sRGB profile that was probably erroneously applied to the image,
	// we can't just apply sRGB with the icc_transform because sRGB isn't
	// an appropriate profile for the icc_transform command, so we have to call
	// vipsthumbnail instead which does some magick behind the scenes to properly
	// convert between the profiles.
	if channels == "gray" && (iccProfileName == "" || iccProfileName == "sRGB Profile") {
		log.Printf("Fixing gray image %s with profile [%s]", c.NoalphaFile, iccProfileName)
		err = vips.FixGray(fmt.Sprintf("%s", c.NoalphaFile), c.GrayFixedFile)
		if err != nil {
			return fmt.Errorf("Agent#toPyramidTIFF FixGray failed - %v", err)
		}
		newProfile = true
	} else {
		c.GrayFixedFile = c.NoalphaFile
	}

	if iccProfileName == "" {
		log.Printf("WARNING icc profile not available for image %s - profile won't be converted\n", c.GrayFixedFile)
	}

	// Some notes:
	// - If no ICC profile is embedded, browsers will usually assume the image is in sRGB.
	// - When ICC profile description string is "sRGB.icc", vips complained it is not complained that
	//   it is not compatible with the the destination profile (sRGB IEC61966-2.1).
	if !newProfile && iccProfileName != "" && !strings.HasPrefix(strings.ToLower(iccProfileName), "srgb") {
		fmt.Printf("ICC transform %s -> %s (%s)\n", c.GrayFixedFile, c.ProfileFixedFile, targetICCProfile)
		err = vips.ICCTransform(fmt.Sprintf("%s[0]", c.GrayFixedFile), c.ProfileFixedFile, targetICCProfile)
		if err != nil {
			return fmt.Errorf("Agent#toPyramidTIFF ICCTransform failed - %v", err)
		}
	} else {
		c.ProfileFixedFile = c.GrayFixedFile
	}

	err = a.createPyramid(c, c.ProfileFixedFile, c.Input.OutFile)
	if err != nil {
		return fmt.Errorf("Agent#toPyramidTIFF createPyramid failed - %v", err)
	}
	return nil
}

func (a *Agent) createPyramid(c *context.Context, inFile, outFile string) (err error) {
	var w, h uint

	if w, h, err = a.initialResize(c, inFile); err != nil {
		return fmt.Errorf("Agent#createTIFF initialResize failed - %v", err)
	}
	c.Output.OutputWidth = w
	c.Output.OutputHeight = h

	if err = a.createSubImages(c, w, h); err != nil {
		return fmt.Errorf("Agent#createTIFF createSubImages failed - %v", err)
	}
	if err = a.combineSubImages(c); err != nil {
		return fmt.Errorf("Agent#createTIFF combineImages failed - %v", err)
	}
	return nil
}

// Prepare the top-level image for the pyramidal TIFF.
func (a *Agent) initialResize(c *context.Context, inFile string) (w, h uint, err error) {
	w, h = c.InitialWH()
	fmt.Printf("initial w: %d, h: %d\n", w, h)
	top := fmt.Sprintf("%s_0.tif", c.TmpFilePrefix)
	inFile0 := fmt.Sprintf("%s[0]", inFile)

	if w == c.Width {
		// Use the original size since it isn't bigger than maxSize.
		log.Printf("Copying %s to %s\n", inFile, top)
		if _, err = util.CopyFile(inFile, top); err != nil {
			log.Printf("ERROR Context#toPyramid copyFile failed - %v\n", err)
			return w, h, err
		}
	} else {
		// Resize original to maxSize.
		err = vips.Resize(inFile0, top, w, h)
		if err != nil {
			log.Printf("ERROR Context#toPyramid - %v\n", err)
		}
	}
	return w, h, err
}

// Create sub-images for the pyramid.
func (a *Agent) createSubImages(c *context.Context, w, h uint) (err error) {
	depth := 1

	for w, h, depth = w/2, h/2, 1; w > 0 && h > 0 && (w > 127 || h > 127); depth++ {
		inFile := fmt.Sprintf("%s_%d.tif", c.TmpFilePrefix, depth-1)
		outFile := fmt.Sprintf("%s_%d.tif", c.TmpFilePrefix, depth)

		if err = vips.Resize(inFile, outFile, w, h); err != nil {
			return err
		}
		w /= 2
		h /= 2
	}
	return err
}

func (a *Agent) combineSubImages(c *context.Context) error {
	inFiles, err := filepath.Glob(fmt.Sprintf("%s_*.tif", c.TmpFilePrefix))

	compression := c.CompressionOption()
	if c.BitDepth > 8 {
		compression = "" // no compression for depth 16 images (jpeg can't handle 16 bit)
		log.Printf("WARNING: JPEG can't handle 16 bit images, so no compression applied\n")
	}

	err = tiff.BuildPyramid(inFiles, c.Input.OutFile, map[string]string{
		"c": compression,
	})

	if err != nil {
		return fmt.Errorf("Agent#combineSubImages failed to build pyramid - %v", err)
	}
	return nil
}

func (a *Agent) validateChannels(channels string) bool {
	switch channels {
	case "srgb", "gray", "cmyk", "srgba":
		return true
	default:
		return false
	}
}

func (a *Agent) mkdirp(d string) error {
	log.Printf("pyramid.agent.Agent#mkdirp making sure directory %s exists", d)
	err := os.MkdirAll(d, 0700)
	if err != nil {
		return fmt.Errorf("pyramid.Agent#mkdirp failed to create directory %s - %v", d, err)
	}
	return nil
}
