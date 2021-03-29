package imagemagick

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentify(t *testing.T) {
	t.Run("GetInfo", func(t *testing.T) {
		imageFormat, channels, depth, profile, err := GetInfo("../../test/resources/images/ag-obj-286-0033-pub.jpg")
		if err != nil {
			t.Errorf("GetInfo - %v", err)
		}
		assert.Equal(t, "JPEG", imageFormat, "%[m]")
		assert.Equal(t, "srgb", channels, "%[channels]")
		assert.Equal(t, "8", depth, "%[bit-depth")
		assert.Equal(t, "Adobe RGB (1998)", profile, "%[profile:icc]")
	})
}
