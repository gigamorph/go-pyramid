package config

import (
	"fmt"
	"os"
)

// TempDir is the directory which holds temporary image files
var TempDir = getEnv("TEMP_DIR", fmt.Sprintf("%s/%s", os.TempDir(), "go-pyramid"))

/************************************************************
 * BEGIN Command-line paths                                 *
 * Where to find external programs to run from the shell.   *
 * These need not be set if external programs are not used. *
 ************************************************************/

// Identify is the path to identify (of ImageMagick).
var Identify = getEnv("IDENTIFY", "/usr/local/bin/identify")

// TIFFCopy is the path to tiffcp
var TIFFCopy = getEnv("TIFFCP", "/usr/local/bin/tiaffcp")

// VIPS is the path to vips.
var VIPS = getEnv("VIPS", "/opt/vips/bin/vips")

// VIPSHeader is the path to vipsheader
var VIPSHeader = getEnv("VIPS_HEADER", "/opt/vips/bin/vipsheader")

// VIPSThumbnail is the path to vipsthumbnail
var VIPSThumbnail = getEnv("VIPS_THUMBNAIL", "/opt/vips/bin/vipsthumbnail")

// ICCProfile is the path to sRGBProfile.icc
var ICCProfile = getEnv("ICC_PROFILE", "/opt/docker-shared/image-conversion-test/sRGBProfile.icc")

/**************************
 * END Command-line paths *
 **************************/

// Returns the value of the environment variable named name.
// If undefined or empty, it returns defaultValue instead.
func getEnv(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		value = defaultValue
	}
	return value
}
