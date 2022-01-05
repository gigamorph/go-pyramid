package combined

import (
	"log"

	"github.com/gigamorph/go-pyramid/config"
	"github.com/gigamorph/go-pyramid/shellcmds/vips"
	"github.com/gigamorph/go-pyramid/util"
)

func GrayToSRGB(inFile, outFile string) error {
	var w, h uint
	var err error

	if w, err = vips.Width(inFile); err != nil {
		return err
	}
	if h, err = vips.Height(inFile); err != nil {
		return err
	}
	log.Printf("width: %d, height: %d", w, h)

	args := []string{
		"-colorspace",
		"srgb",
		"-type",
		"truecolor",
		inFile,
		outFile,
	}
	_, err = util.Exec(config.Convert, args)
	return err
}
