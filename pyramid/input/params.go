package input

// Params holds user-provided parameters.
type Params struct {
	InFile           string
	OutFile          string
	MaxSize          uint   // max outfile size (long-edge)
	Compression      string // compression method ("jpeg", "lzw", "")
	Quality          int    // JPEG quality (1-100)
	TargetICCProfile string // file path of the profile
	TempDir          string // path of directory where temporary files will be stored

	// Path of directory where ImageMagick temporary files will be stored.
	// If nil, default will be used.
	IMTempDir *string

	DeleteTemp bool // delete temp dir after conversion is done
}
