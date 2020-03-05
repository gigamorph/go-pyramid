package imagemagick

import (
	"fmt"

	"github.com/gigamorph/go-pyramid/util"
	"gopkg.in/gographics/imagick.v3/imagick"
)

// IM encapsulates the ImageMagick library
type IM struct {
	mw *imagick.MagickWand
}

// GetIM returns a single instance of IM
func GetIM() *IM {
	if instance == nil {
		imagick.Initialize()
		instance = new()
	}
	return instance
}

var instance *IM

func new() *IM {
	im := IM{
		mw: imagick.NewMagickWand(),
	}
	return &im
}

// Finalize cleans up after the ImageMagick library
func (im *IM) Finalize() {
	im.mw.Destroy()
	imagick.Terminate()
}

// ReadImage reads image into memory for further operation
func (im *IM) ReadImage(fpath string) error {
	err := im.mw.ReadImage(fpath)
	if err != nil {
		return fmt.Errorf("IM#ReadImage failed - %v", err)
	}
	return nil
}

// GetImageWidth returns the width of the image
func (im *IM) GetImageWidth() uint {
	return im.mw.GetImageWidth()
}

// GetImageHeight returns the height of the image
func (im *IM) GetImageHeight() uint {
	return im.mw.GetImageHeight()
}

// GetImageColorspace returns the colorspace ID of the image
func (im *IM) GetImageColorspace() imagick.ColorspaceType {
	return im.mw.GetImageColorspace()
}

// GetImageAlphaChannel returns true if the image has an alpha channel
func (im *IM) GetImageAlphaChannel() bool {
	return im.mw.GetImageAlphaChannel()
}

// GetImageFormat returns format of the image in string; e.g. "TIFF"
func (im *IM) GetImageFormat() string {
	return im.mw.GetImageFormat()
}

// GetICCProfileName returns ICC profile name embedded in the image.
// If not available, it returns an empty string "".
func (im *IM) GetICCProfileName() string {
	iccProfile := im.mw.GetImageProfile("ICC")
	iccProfileName := util.GetICCProfileDesc([]byte(iccProfile))
	return iccProfileName
}
