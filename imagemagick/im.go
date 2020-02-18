package imagemagick

import (
	"fmt"

	"github.com/gigamorph/go-pyramid/util"
	"gopkg.in/gographics/imagick.v3/imagick"
)

type IM struct {
	mw *imagick.MagickWand
}

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

func (im *IM) Finalize() {
	im.mw.Destroy()
	imagick.Terminate()
}

func (im *IM) ReadImage(fpath string) error {
	err := im.mw.ReadImage(fpath)
	if err != nil {
		return fmt.Errorf("IM#ReadImage failed - %v", err)
	}
	return nil
}

func (im *IM) GetImageWidth() uint {
	return im.mw.GetImageWidth()
}

func (im *IM) GetImageHeight() uint {
	return im.mw.GetImageHeight()
}

func (im *IM) GetImageColorspace() imagick.ColorspaceType {
	return im.mw.GetImageColorspace()
}

func (im *IM) GetImageAlphaChannel() bool {
	return im.mw.GetImageAlphaChannel()
}

func (im *IM) GetImageFormat() string {
	return im.mw.GetImageFormat()
}

func (im *IM) GetICCProfileName() string {
	iccProfile := im.mw.GetImageProfile("ICC")
	iccProfileName := util.GetICCProfileDesc([]byte(iccProfile))
	return iccProfileName
}
