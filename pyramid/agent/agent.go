package agent

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/gigamorph/go-pyramid/config"
	"github.com/gigamorph/go-pyramid/imagemagick"
	"github.com/gigamorph/go-pyramid/pyramid/context"
	"github.com/gigamorph/go-pyramid/pyramid/input"
	"github.com/gigamorph/go-pyramid/pyramid/output"
	"github.com/gigamorph/go-pyramid/shellcmds"
	"github.com/gigamorph/go-pyramid/util"
	"github.com/gigamorph/go-pyramid/vips"
	"gopkg.in/gographics/imagick.v3/imagick"
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
	im   *imagemagick.IM
	vips *vips.VIPS
}

// New returns a new instance of Agent.
func New(im *imagemagick.IM, vips *vips.VIPS) *Agent {
	agent := Agent{
		im:   im,
		vips: vips,
	}
	return &agent
}

// Convert is the public method to call to actually convert an image.
// p contains input, output, and other information needed for conversion.
func (a *Agent) Convert(p input.Params) (*output.Params, error) {
	im := a.im
	context := context.New(p)
	err := im.ReadImage(p.InFile)
	if err != nil {
		return nil, fmt.Errorf("Agent#Convert failed to read image - %v", err)
	}

	context.Width = im.GetImageWidth()
	context.Height = im.GetImageHeight()
	context.Output.InputWidth = context.Width
	context.Output.InputHeight = context.Height

	err = a.toPyramidTIFF(context)
	if err != nil {
		return nil, fmt.Errorf("Agent#Convert PyramidTIFF failed - %v", err)
	}
	return &context.Output, nil
}

func (a *Agent) toPyramidTIFF(c *context.Context) (err error) {
	var imageFormat, iccProfile string
	im := a.im
	v := a.vips

	p := c.Input
	inFile, outFile := p.InFile, p.OutFile

	targetICCProfile := config.TargetICCProfileIIIF
	if c.Input.TargetICCProfile != "" {
		targetICCProfile = c.Input.TargetICCProfile
	}

	colorspace := im.GetImageColorspace()
	hasAlpha := im.GetImageAlphaChannel()
	imageFormat = im.GetImageFormat()
	log.Printf("colorspace: %d, hasAlpha: %t, imageFormat: %s\n", colorspace, hasAlpha, imageFormat)

	iccProfileName := im.GetICCProfileName()
	log.Printf("Source ICC: [%s]\n", iccProfileName)

	// Check if channels is supported
	if valid := a.validateChannels(colorspace); !valid {
		log.Printf("ERROR Image %s has colorspace %d which is not supported at this time.\n",
			inFile, colorspace)
		return fmt.Errorf("Invalid colorspace %d for %s", colorspace, inFile)
	}

	// Make sure the input is a TIFF file
	if imageFormat != "TIFF" {
		if err = v.ToTiff(inFile, c.TiffFile); err != nil {
			return fmt.Errorf("Agent#toPyramidTIFF ToTiff failed - %v", err)
		}
	} else {
		c.TiffFile = inFile
	}

	// We have to flatten the image to remove the alpha channel / trasparency
	// before proceeding
	// if channels == "srgba" {
	if hasAlpha {
		if err = v.RemoveAlpha(c.TiffFile, c.NoalphaFile); err != nil {
			return fmt.Errorf("Agent#toPyramidTIFF RemoveAlpha failed - %v", err)
		}
	} else {
		c.NoalphaFile = c.TiffFile
	}

	newProfile := false

	// In the case of gray with no embedded color profile or with an embedded
	// sRGB profile that was probably erroneously applied to the image,
	// we can't just apply sRGB with the icc_transform because sRGB isn't
	// an appropriate profile for the icc_transform command, so we have to call
	// vipsthumbnail instead which does some magick behind the scenes to properly
	// convert between the profiles.
	if colorspace == imagick.COLORSPACE_GRAY && (iccProfileName == "" || iccProfileName == "sRGB Profile") {
		log.Printf("Fixing gray image %s with profile %s", inFile, iccProfile)
		if err = v.FixGray(c.NoalphaFile, c.GrayFixedFile, c.Width, targetICCProfile); err != nil {
			return fmt.Errorf("Agent#toPyramidTIFF FixGray failed - %v", err)
		}
		newProfile = true
	} else {
		c.GrayFixedFile = c.NoalphaFile
	}

	if iccProfileName == "" {
		log.Printf("WARNING icc profile not available for image %s - profile won't be converted\n", inFile)
	}

	if !newProfile && iccProfileName != "" {
		fmt.Printf("ICC transform %s -> %s\n", c.GrayFixedFile, c.ProfileFixedFile)
		if err = a.vips.ICCTransformFile(c.GrayFixedFile, c.ProfileFixedFile, targetICCProfile); err != nil {
			return fmt.Errorf("Agent#toPyramidTIFF ICCTransform failed - %v", err)
		}
	} else {
		c.ProfileFixedFile = c.GrayFixedFile
	}

	err = a.createPyramid(c, c.ProfileFixedFile, outFile)
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
	v := a.vips
	w, h = c.InitialWH()
	fmt.Printf("initial w: %d, h: %d\n", w, h)
	top := fmt.Sprintf("%s_0.tif", c.TmpFilePrefix)

	if w == c.Width {
		// Use the original size since it isn't bigger than maxSize.
		log.Printf("Copying %s to %s\n", inFile, top)
		if _, err = util.CopyFile(inFile, top); err != nil {
			log.Printf("ERROR Context#toPyramid copyFile failed - %v\n", err)
			return w, h, err
		}
	} else {
		// Resize original to maxSize.
		err = v.Resize(inFile, top, w, h)
		if err != nil {
			log.Printf("ERROR Context#toPyramid - %v\n", err)
		}
	}
	return w, h, err
}

// Create sub-images for the pyramid.
func (a *Agent) createSubImages(c *context.Context, w, h uint) (err error) {
	v := a.vips
	depth := 1

	fmt.Printf("SUB %d %d\n", w, h)

	for w, h, depth = w/2, h/2, 1; w > 0 && h > 0 && (w > 127 || h > 127); depth++ {
		fmt.Printf("SUB %d %d\n", w, h)
		inFile := fmt.Sprintf("%s_%d.tif", c.TmpFilePrefix, depth-1)
		outFile := fmt.Sprintf("%s_%d.tif", c.TmpFilePrefix, depth)

		if err = v.Resize(inFile, outFile, w, h); err != nil {
			return err
		}
		w /= 2
		h /= 2
	}
	return err
}

func (a *Agent) combineSubImages(c *context.Context) error {
	inFiles, err := filepath.Glob(fmt.Sprintf("%s_*.tif", c.TmpFilePrefix))

	err = shellcmds.BuildPyramid(inFiles, c.Input.OutFile, map[string]string{
		"c": c.CompressionOption(),
	})

	if err != nil {
		return fmt.Errorf("Agent#combineSubImages failed to build pyramid - %v", err)
	}
	return nil
}

func (a *Agent) validateChannels(channels imagick.ColorspaceType) bool {
	switch channels {
	case imagick.COLORSPACE_SRGB, imagick.COLORSPACE_GRAY, imagick.COLORSPACE_CMYK:
		return true
	default:
		return false
	}
}
