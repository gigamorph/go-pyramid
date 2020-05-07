package tiff

import (
	"fmt"

	"github.com/gigamorph/go-pyramid/config"
	"github.com/gigamorph/go-pyramid/util"
)

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
