package context

import (
	"fmt"
	"math"
	"os"
	"path"
	"strings"

	"github.com/gigamorph/go-pyramid/pyramid/input"
	"github.com/gigamorph/go-pyramid/pyramid/output"
)

// Context holds inforamtion needed to perform conversion.
type Context struct {
	Input            input.Params
	Output           output.Params
	TmpFilePrefix    string
	TiffFile         string
	NoalphaFile      string
	GrayFixedFile    string
	ProfileFixedFile string
	Width            uint // original width
	Height           uint // original height
	BitDepth         uint // original bit depth, e.g. 8, 16
}

// New returns a new instance of Context.
func New(p input.Params) *Context {
	c := Context{}

	base := path.Base(p.InFile)
	ext := path.Ext(base)
	name := strings.TrimSuffix(base, ext)

	c.Input = p
	if c.Input.TempDir == "" {
		c.Input.TempDir = fmt.Sprintf("%s/cds2-scaler", os.TempDir())
	}
	if c.Input.Compression == "" {
		c.Input.Compression = "jpeg"
		c.Input.Quality = 90
	}
	c.TmpFilePrefix = fmt.Sprintf("%s/%s", c.Input.TempDir, name)
	c.TiffFile = fmt.Sprintf("%s.tif", c.TmpFilePrefix)
	c.NoalphaFile = fmt.Sprintf("%s.noalpha.tif", c.TmpFilePrefix)
	c.GrayFixedFile = fmt.Sprintf("%s.grayfixed.tif", c.TmpFilePrefix)
	c.ProfileFixedFile = fmt.Sprintf("%s.profilefixed.tif", c.TmpFilePrefix)
	return &c
}

// InitialWH calculates the size of the biggest tile in the output pyramidal TIFF
func (c *Context) InitialWH() (uint, uint) {
	w0, h0 := c.Width, c.Height // original dimensions
	w, h := w0, h0              // new starting point
	max := c.Input.MaxSize

	if max != 0 {
		if w < h {
			if h > max {
				h = max
				w = uint(math.Round(float64(w0*h) / float64(h0)))
			}
		} else { // w >= h
			if w > max {
				w = max
				h = uint(math.Round(float64(h0*w) / float64(w0)))
			}
		}
	}
	return w, h
}

func (c *Context) CompressionOption() string {
	switch c.Input.Compression {
	case "jpeg":
		quality := c.Input.Quality
		return fmt.Sprintf("jpeg:%d", quality)
	case "lzw":
		return "lzw"
	default:
		return ""
	}
}
