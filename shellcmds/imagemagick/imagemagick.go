package imagemagick

import (
	"fmt"

	"github.com/gigamorph/go-pyramid/config"
	"github.com/gigamorph/go-pyramid/util"
)

// ImageFormat returns the "magick" value, e.g. "TIFF", "JPEG"
func ImageFormat(fpath string) (string, error) {
	var out string

	args := []string{
		"-format", "%[m]",
		fmt.Sprintf("%s[0]", fpath),
	}

	out, err := util.Exec(config.Identify, args)
	if err != nil {
		return "", err
	}
	return out, nil
}

// Channels returns the channels string acquired from the image file by ImageMagick/identify.
func Channels(fpath string) (channels string, err error) {
	var out string

	args := []string{
		"-format", "%[channels]",
		fmt.Sprintf("%s[0]", fpath),
	}

	if out, err = util.Exec(config.Identify, args); err != nil {
		return "", err
	}
	return out, err
}

// ICCProfile returns the ICC profile identifier string acquired from
// the image by ImageMagic/identify.
func ICCProfile(fpath string) (iccProfile string, err error) {
	var out string

	args := []string{
		"-format", "%[profile:icc]",
		fmt.Sprintf("%s[0]", fpath),
	}

	if out, err = util.Exec(config.Identify, args); err != nil {
		return "", err
	}
	return out, err
}
