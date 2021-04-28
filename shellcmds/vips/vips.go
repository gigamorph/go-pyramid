package vips

import (
	"fmt"
	"log"
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

// RemoveAlpha strippes the alpha channel from inFile.
func RemoveAlpha(inFile, outFile string) error {
	args := []string{
		"im_extract_bands",
		inFile,
		outFile,
		"0",
		"3",
	}
	_, err := util.Exec(config.VIPS, args)
	return err
}

// FixGray fixes some issues with "gray" images.
//
// In the case of gray with no embedded color profile or with an embedded
// sRGB profile that was probably erroneously applied to the image,
// we can't just apply sRGB with the icc_transform because sRGB isn't
// an appropriate profile for the icc_transform command so we have to
// call vipsthumbnail instead which does some magick behind the scenes
// to properly convert between the profiles.
func FixGray(inFile, outFile string) error {
	var w, h uint
	var err error

	if w, err = Width(inFile); err != nil {
		return err
	}
	if h, err = Height(inFile); err != nil {
		return err
	}
	log.Printf("width: %d, height: %d", w, h)

	args := []string{
		inFile,
		fmt.Sprintf("--eprofile=%s", config.TargetICCProfileIIIF),
		"--size", fmt.Sprintf("%dx%d", w, h),
		"--intent", "relative",
		"-o", fmt.Sprintf("%s[compression=none,strip]", outFile),
	}
	_, err = util.Exec(config.VIPSThumbnail, args)
	return err
}

// ICCTransform changes the color profile.
func ICCTransform(inFile, outFile, iccProfile string) error {
	args := []string{
		"icc_transform",
		inFile,
		//fmt.Sprintf("%s[compression=none,strip]", outFile),
		fmt.Sprintf("%s[compression=none]", outFile),
		iccProfile,
		"--embedded",
		"--input-profile", config.TargetICCProfileIIIF,
		"--intent", "relative",
	}
	_, err := util.Exec(config.VIPS, args)
	return err
}

// Resize the image.
func Resize(inFile, outFile string, width, height uint) error {
	args := []string{
		inFile,
		"--size", fmt.Sprintf("%dx%d!", width, height),
		"-o", outFile,
	}
	_, err := util.Exec(config.VIPSThumbnail, args)
	return err
}

// ResizeBoundedNoExpand resizes the image to fit the bounding box of width x height
// with aspect ratio preserved, but does not resize it if the source image
// is smaller
func ResizeBoundedNoExpand(inFile, outFile string, width, height uint) error {
	args := []string{
		inFile,
		"--size", fmt.Sprintf("%dx%d>", width, height),
		"-o", outFile,
	}
	_, err := util.Exec(config.VIPSThumbnail, args)
	return err
}

// ToTiff converts inFile to TIFF format.
func ToTiff(inFile, outFile string) error {
	args := []string{
		"tiffsave",
		inFile,
		outFile,
	}
	_, err := util.Exec(config.VIPS, args)
	return err
}
