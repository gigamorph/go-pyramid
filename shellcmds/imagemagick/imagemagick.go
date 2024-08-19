package imagemagick

import (
	"fmt"
	"log"
	"strings"

	"github.com/gigamorph/go-pyramid/config"
	"github.com/gigamorph/go-pyramid/shellcmds/vips"
	"github.com/gigamorph/go-pyramid/util"
)

func tempDirArg(tempDir string) string {
	return fmt.Sprintf("registry:temporary-path=%s", tempDir)
}

// ImageFormat returns the "magick" value, e.g. "TIFF", "JPEG"
func ImageFormat(fpath string, tempDir *string) (string, error) {
	var out string

	args := make([]string, 0, 5)
	args = append(args, "-format", "%[m]")
	if tempDir != nil {
		args = append(args, "-define", tempDirArg(*tempDir))
	}
	args = append(args, fmt.Sprintf("%s[0]", fpath))

	out, err := util.Exec(config.Identify, args)
	if err != nil {
		return "", err
	}
	return out, nil
}

// Channels returns the channels string acquired from the image file by ImageMagick/identify.
func Channels(fpath string, tempDir *string) (channels string, err error) {
	var out string

	args := make([]string, 0, 5)
	args = append(args, "-format", "%[channels]")
	if tempDir != nil {
		args = append(args, "-define", tempDirArg(*tempDir))
	}
	args = append(args, fmt.Sprintf("%s[0]", fpath))

	if out, err = util.Exec(config.Identify, args); err != nil {
		return "", err
	}
	return out, err
}

// ICCProfile returns the ICC profile identifier string acquired from
// the image by ImageMagic/identify.
func ICCProfile(fpath string, tempDir *string) (iccProfile string, err error) {
	var out string

	args := make([]string, 0, 5)
	args = append(args, "-format", "%[profile:icc]")
	if tempDir != nil {
		args = append(args, "-define", tempDirArg(*tempDir))
	}
	args = append(args, fmt.Sprintf("%s[0]", fpath))

	if out, err = util.Exec(config.Identify, args); err != nil {
		return "", err
	}
	return out, err
}

// GetInfo returns multiple information from identify.
// Running identify for those separately is very costly for large images.
func GetInfo(fpath string, tempDir *string) (string, string, string, string, error) {
	args := make([]string, 0, 5)
	args = append(args, "-format", "%[m]|%[channels]|%[bit-depth]|%[profile:icc]")
	if tempDir != nil {
		args = append(args, "-define", tempDirArg(*tempDir))
	}
	args = append(args, fmt.Sprintf("%s[0]", fpath))

	out, err := util.Exec(config.Identify, args)
	if err != nil {
		return "", "", "", "", fmt.Errorf("imagemagick.GetInfo failed - %v", err)
	}
	values := strings.Split(out, "|")
	return values[0], values[1], values[2], values[3], err
}

func GrayToSRGB(inFile, outFile string, tempDir *string) error {
	var w, h uint
	var err error

	if w, err = vips.Width(inFile); err != nil {
		return err
	}
	if h, err = vips.Height(inFile); err != nil {
		return err
	}
	log.Printf("width: %d, height: %d", w, h)

	args := make([]string, 0, 5)
	args = append(args, inFile,
		fmt.Sprintf("--eprofile=%s", config.TargetICCProfileIIIF),
		"--size", fmt.Sprintf("%dx%d", w, h),
		"--intent", "relative",
		"-o", fmt.Sprintf("%s[compression=none,strip]", outFile),
	)
	if tempDir != nil {
		args = append(args, "-define", tempDirArg(*tempDir))
	}

	_, err = util.Exec(config.VIPSThumbnail, args)
	return err
}
