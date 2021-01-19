package util

import (
	"fmt"
	"math"
)

// ThumbnailSizeByLongSide calculates the size of the thumbnail
// when maximum size of the longer side is provided.
// If the original image is smaller than the imposed limit,
// the original dimensions are returned.
func ThumbnailSizeByLongSide(width, height, maxLong uint) (uint, uint, error) {
	if width < 1 || height < 1 || maxLong < 1 {
		return 0, 0, fmt.Errorf("Invalid input %d, %d, %d", width, height, maxLong)
	}
	// Do not resize if the original image is smaller than maxLong
	if width <= maxLong && height <= maxLong {
		return width, height, nil
	}
	fw, fh, fmax := float64(width), float64(height), float64(maxLong)
	longer, shorter := fw, fh
	if width < height {
		longer, shorter = fh, fw
	}
	scale := longer / fmax
	if width < height {
		return uint(math.Round(shorter / scale)), uint(math.Round(longer / scale)), nil
	}
	return uint(math.Round(longer / scale)), uint(math.Round(shorter / scale)), nil
}
