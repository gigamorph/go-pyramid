// Usage:
// go run pyramid.go [options] <infile> <outfile>
// options: -m, -c, -q, -p, -t (see main function below)
package main

import (
	"flag"
	"log"

	"github.com/gigamorph/go-pyramid/pyramid/agent"
	"github.com/gigamorph/go-pyramid/pyramid/input"
)

func main() {
	maxSizePtr := flag.Uint("m", 0, "max size")
	compressionPtr := flag.String("c", "", "compression method")
	qualityPtr := flag.Int("q", 90, "jpeg quality (1-100)")
	targetProfilePtr := flag.String("p", "test/resources/sRGBProfile.icc", "ICC profile of target file")
	tempDirPtr := flag.String("t", "/tmp/go-pyramid", "path to temp dir")
	flag.Parse()

	args := flag.Args()

	inFile := args[0]
	outFile := args[1]

	log.Printf("BEGIN processing image file %s\n", inFile)
	log.Printf("maxSize: %d\n", *maxSizePtr)

	params := input.Params{
		InFile:           inFile,
		OutFile:          outFile,
		MaxSize:          *maxSizePtr,
		Compression:      *compressionPtr,
		Quality:          *qualityPtr,
		TargetICCProfile: *targetProfilePtr,
		TempDir:          *tempDirPtr,
	}

	ag := agent.New()

	_, err := ag.Convert(params)
	if err != nil {
		log.Printf("ERROR main agent.Convert failed - %v", err)
	}
}
