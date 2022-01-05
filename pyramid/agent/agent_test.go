package agent

import (
	"fmt"
	"testing"

	"github.com/gigamorph/go-pyramid/pyramid/input"
)

func TestConvert(t *testing.T) {
	a := New()

	tempDir := "/tmp/go-pyramid-test"
	iccProfile := fromRoot("test/resources/sRGBProfile.icc")

	t.Run("RemoveAlpha", func(t *testing.T) {
		_, err := a.Convert(input.Params{
			InFile:           fromRoot("test/resources/images/ag-obj-286-0033-pub.alpha.tif"),
			OutFile:          fromRoot("tmp/ag.alpha.tif"),
			MaxSize:          0,
			Compression:      "jpeg",
			Quality:          90,
			TargetICCProfile: iccProfile,
			TempDir:          tempDir,
			DeleteTemp:       true,
		})
		if err != nil {
			t.Errorf("RemoveAlpha - %v", err)
		}
	})

	t.Run("ConvertJPEG", func(t *testing.T) {
		_, err := a.Convert(input.Params{
			InFile:           fromRoot("test/resources/images/ag-obj-286-0033-pub.jpg"),
			OutFile:          fromRoot("tmp/ag.jpg.tif"),
			MaxSize:          0,
			Compression:      "jpeg",
			Quality:          90,
			TargetICCProfile: iccProfile,
			TempDir:          tempDir,
		})
		if err != nil {
			t.Errorf("ConvertJPEG - %v", err)
		}
	})

	t.Run("GrayscaleWithAdobeRGB", func(t *testing.T) {
		_, err := a.Convert(input.Params{
			InFile:           fromRoot("test/resources/images/grayscale-with-adobe-rgb-1998.tif"),
			OutFile:          fromRoot("tmp/gray-adobe-rgb.tif"),
			MaxSize:          0,
			Compression:      "jpeg",
			Quality:          90,
			TargetICCProfile: iccProfile,
			TempDir:          tempDir,
			DeleteTemp:       true,
		})
		if err != nil {
			t.Errorf("GrayscaleWithAdobeRGB - %v", err)
		}
	})
}

func fromRoot(relPath string) string {
	return fmt.Sprintf("../../%s", relPath)
}
