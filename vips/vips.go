package vips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

import (
	"fmt"

	"github.com/davidbyttow/govips/pkg/vips"
	"github.com/gigamorph/go-pyramid/config"
)

type VIPS struct {
}

func GetVIPS() *VIPS {
	if instance == nil {
		vips.Startup(nil)
		instance = new()
	}
	return instance
}

var instance *VIPS

func new() *VIPS {
	v := VIPS{}
	return &v
}

func (v *VIPS) Finalize() {
	vips.ShutdownThread()
	vips.Shutdown()
}

// ToTiff converts inFile to TIFF format.
func (v *VIPS) ToTiff(inFile, outFile string) error {
	image, err := vips.NewImageFromFile(inFile)
	if err != nil {
		return fmt.Errorf("VIPS#ToTiff failed to read image - %v", err)
	}
	err = vips.Tiffsave(image.Image(), outFile)
	if err != nil {
		return fmt.Errorf("VIPS#ToTiff failed to save tiff - %v", err)
	}
	return nil
}

// Resize the image.
func (v *VIPS) Resize(inFile, outFile string, width, height uint) error {
	image, err := vips.Thumbnail(inFile, int(width),
		vips.InputInt("height", int(height)),
	)
	if err != nil {
		return fmt.Errorf("VIPS#Resize failed to create thumbnail - %v", err)
	}

	err = vips.Tiffsave(image, outFile)
	if err != nil {
		return fmt.Errorf("VIPS#Resize failed to save tiff - %v", err)
	}
	return nil
}

// RemoveAlpha strippes the alpha channel from inFile.
func (v *VIPS) RemoveAlpha(inFile, outFile string) error {
	image, err := vips.NewImageFromFile(inFile)
	if err != nil {
		return fmt.Errorf("VIPS#ToTiff failed to read image - %v", err)
	}
	outImage, err := vips.ExtractBand(image.Image(), 0, vips.InputInt("n", 3))
	if err != nil {
		return fmt.Errorf("VIPS#ToTiff failed to extract bands - %v", err)
	}
	err = vips.Tiffsave(outImage, outFile)
	if err != nil {
		return fmt.Errorf("VIPS#RemoveAlpha failed to save image - %v", err)
	}
	return nil
}

func (v *VIPS) ICCTransformFile(inFile string, outFile string, outProfilePath string) error {
	image, err := vips.NewImageFromFile(inFile)
	if err != nil {
		return fmt.Errorf("VIPS#ICCTransformFile failed to read image - %v", err)
	}

	// Default intent is VIPS_INTENT_RELATIVE:
	// see https://libvips.github.io/libvips/API/current/libvips-resample.html
	outImage, err := vips.IccTransform(image.Image(), outProfilePath)
	if err != nil {
		return fmt.Errorf("VIPS#ICCTransformFile failed to transform - %v", err)
	}
	err = vips.Tiffsave(outImage, outFile)
	if err != nil {
		return fmt.Errorf("VIPS#ICCTransformFile failed to save image - %v", err)
	}
	return nil
}

// FixGray fixes some issues with "gray" images.
//
// In the case of gray with no embedded color profile or with an embedded
// sRGB profile that was probably erroneously applied to the image,
// we can't just apply sRGB with the icc_transform because sRGB isn't
// an appropriate profile for the icc_transform command so we have to
// call vipsthumbnail instead which does some magick behind the scenes
// to properly convert between the profiles.
func (v *VIPS) FixGray(inFile, outFile string, width uint) error {
	image, err := vips.Thumbnail(inFile, int(width),
		vips.InputString("export_profile", config.ICCProfile),
		vips.InputString("intent", "perceptual"),
	)
	if err != nil {
		return fmt.Errorf("FixGray failed to create thumbnail - %v", err)
	}

	err = vips.Tiffsave(image, outFile)
	if err != nil {
		return fmt.Errorf("VIPS#FixGray failed to save tiff - %v", err)
	}
	return nil
}
