package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThumbnailSizeByLongSide(t *testing.T) {
	w, h, err := ThumbnailSizeByLongSide(800, 600, 480)
	assert.Equal(t, nil, err, "Width is longer - should cause no error")
	assert.Equal(t, uint(480), w, "Width is longer - correct width")
	assert.Equal(t, uint(360), h, "Width is longer - correct height")

	w, h, err = ThumbnailSizeByLongSide(600, 800, 480)
	assert.Equal(t, nil, err, "Height is longer - should cause no error")
	assert.Equal(t, uint(360), w, "Height is longer - correct width")
	assert.Equal(t, uint(480), h, "Height is longer - correct height")

	w, h, err = ThumbnailSizeByLongSide(800, 800, 480)
	assert.Equal(t, nil, err, "Square - should cause no error")
	assert.Equal(t, uint(480), w, "Square - correct width")
	assert.Equal(t, uint(480), h, "Square - correct height")

	w, h, err = ThumbnailSizeByLongSide(300, 200, 480)
	assert.Equal(t, nil, err, "Smaller than limit - should cause no error")
	assert.Equal(t, uint(300), w, "Smaller than limit - correct width")
	assert.Equal(t, uint(200), h, "Smaller than limit - correct height")

	w, h, err = ThumbnailSizeByLongSide(0, 800, 480)
	assert.NotEqual(t, nil, err, "width = 0 - should cause error")

	w, h, err = ThumbnailSizeByLongSide(800, 0, 480)
	assert.NotEqual(t, nil, err, "height = 0 - should cause error")

	w, h, err = ThumbnailSizeByLongSide(800, 600, 0)
	assert.NotEqual(t, nil, err, "maxLong = 0 - should cause error")
}
