package shellcmds

import (
	"fmt"
	"log"

	"github.com/gigamorph/go-pyramid/config"
	"github.com/gigamorph/go-pyramid/util"
)

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
		"--intent", "perceptual",
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
		config.TargetICCProfileIIIF,
		"--embedded",
		"--input-profile", iccProfile,
		"--intent", "perceptual",
	}
	_, err := util.Exec(config.VIPS, args)
	return err
}

// Resize the image.
func Resize(inFile, outFile string, width, height uint) error {
	args := []string{
		inFile,
		"--size", fmt.Sprintf("%dx%d!", width, height),
		"-o", fmt.Sprintf("%s[compression=none]", outFile),
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

// BuildPyramid contcatenates tiles into one pyramid TIFF
func BuildPyramid(inFiles []string, outFile string, options map[string]string) (err error) {
	args := make([]string, 0, 32)

	// c: compression. e.g.) "jpeg:90"
	if c := options["c"]; c != "" {
		args = append(args, "-c", c)
	}

	args = append(args,
		"-t",        // output to tiles
		"-w", "256", // tile width
		"-l", "256", // tile length
	)
	args = append(args, inFiles...)
	args = append(args, outFile)

	_, err = util.Exec(config.TIFFCopy, args)
	if err != nil {
		return fmt.Errorf("shellcmds.BuildPyramid util.Exec failed - %v", err)
	}
	return nil
}
