package context

import (
	"fmt"
	"math"
	"path"
	"strings"

	"github.com/gigamorph/go-pyramid/config"
	"github.com/gigamorph/go-pyramid/pyramid/input"
	"github.com/gigamorph/go-pyramid/pyramid/output"
)

// Context holds inforamtion needed to perform conversion.
type Context struct {
	Input            input.Params
	Output           output.Params
	TiffFile         string
	NoalphaFile      string
	GrayFixedFile    string
	ProfileFixedFile string
	TmpFilePrefix    string
	Width            uint // original width
	Height           uint // original height
}

// New returns a new instance of Context.
func New(p input.Params) *Context {
	c := Context{}

	base := path.Base(p.InFile)
	ext := path.Ext(base)
	name := strings.TrimSuffix(base, ext)

	c.Input = p
	c.TiffFile = fmt.Sprintf("%s/%s.tif", config.TempDir, name)
	c.NoalphaFile = fmt.Sprintf("%s/%s.noalpha.tif", config.TempDir, name)
	c.GrayFixedFile = fmt.Sprintf("%s/%s.grayfixed.tif", config.TempDir, name)
	c.ProfileFixedFile = fmt.Sprintf("%s/%s.profilefixed.tif", config.TempDir, name)
	c.TmpFilePrefix = fmt.Sprintf("%s/%s", config.TempDir, name)

	return &c
}

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
