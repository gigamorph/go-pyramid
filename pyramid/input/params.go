package input

// Params holds user-provided parameters.
type Params struct {
	InFile           string
	OutFile          string
	MaxSize          uint   // max outfile size (long-edge)
	Compression      string // compression method ("jpeg", "lzw", "")
	Quality          int    // JPEG quality (1-100)
	TargetICCProfile string // file path of the profile
}
