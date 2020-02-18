package main

import (
	"flag"
	"log"

	"github.com/gigamorph/go-pyramid/imagemagick"
	"github.com/gigamorph/go-pyramid/pyramid/agent"
	"github.com/gigamorph/go-pyramid/pyramid/input"
	"github.com/gigamorph/go-pyramid/vips"
)

func main() {
	maxSizePtr := flag.Uint("m", 0, "max size")
	compressionPtr := flag.String("c", "", "compression method")
	qualityPtr := flag.Int("q", 90, "jpeg quality (1-100)")
	flag.Parse()

	args := flag.Args()

	inFile := args[0]
	outFile := args[1]

	log.Printf("BEGIN processing image file %s\n", inFile)
	log.Printf("maxSize: %d\n", *maxSizePtr)

	params := input.Params{
		InFile:      inFile,
		OutFile:     outFile,
		MaxSize:     *maxSizePtr,
		Compression: *compressionPtr,
		Quality:     *qualityPtr,
	}

	im := imagemagick.GetIM()
	defer im.Finalize()

	vips := vips.GetVIPS()
	defer vips.Finalize()

	a := agent.New(im, vips)

	_, err := a.Convert(params)
	if err != nil {
		log.Printf("ERROR main agent.Convert failed - %v", err)
	}
	//log.Printf("Output: %v\n", out)
}
