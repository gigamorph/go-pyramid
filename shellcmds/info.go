package shellcmds

import (
	"fmt"
	"strconv"

	"github.com/gigamorph/go-pyramid/config"
	"github.com/gigamorph/go-pyramid/util"
)

// Width returns the pixel width of the imaage
func Width(fpath string) (w uint, err error) {
	var out string
	var width int64

	args := []string{
		"-f",
		"width",
		fmt.Sprintf("%s[0]", fpath),
	}

	if out, err = util.Exec(config.VIPSHeader, args); err != nil {
		return 0, err
	}

	if width, err = strconv.ParseInt(out, 10, 64); err != nil {
		return 0, err
	}

	return uint(width), err
}

// Height returns the pixel width of the imaage.
func Height(fpath string) (h uint, err error) {
	var out string
	var height int64

	args := []string{
		"-f", "height",
		fmt.Sprintf("%s[0]", fpath),
	}

	if out, err = util.Exec(config.VIPSHeader, args); err != nil {
		return 0, err
	}

	if height, err = strconv.ParseInt(out, 10, 64); err != nil {
		return 0, err
	}

	return uint(height), err
}

// WH returns width and height of the image.
func WH(fpath string) (w, h uint, err error) {
	if w, err = Width(fpath); err != nil {
		return 0, 0, err
	}
	if h, err = Height(fpath); err != nil {
		return 0, 0, err
	}
	return w, h, err
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

// ImageFormat returns the image type of the file e.g., TIFF, JPEG
func ImageFormat(fpath string) (string, error) {
	args := []string{
		"-format", "%[m]",
		fpath,
	}
	out, err := util.Exec(config.Identify, args)
	return out, err
}
