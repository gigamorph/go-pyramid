package config

import (
	"log"
	"os"
)

// TempDir is the directory which holds temporary image files
var TempDir = getEnv("GO_PYRAMID_TEMP_DIR")

/************************************************************
 * BEGIN Command-line paths                                 *
 * Where to find external programs to run from the shell.   *
 * These need not be set if external programs are not used. *
 ************************************************************/

// Identify is the path to identify (of ImageMagick).
var Identify = getEnv("IDENTIFY")

// TIFFCopy is the path to tiffcp
var TIFFCopy = getEnv("TIFFCP")

// VIPS is the path to vips.
var VIPS = getEnv("VIPS")

// VIPSHeader is the path to vipsheader
var VIPSHeader = getEnv("VIPS_HEADER")

// VIPSThumbnail is the path to vipsthumbnail
var VIPSThumbnail = getEnv("VIPS_THUMBNAIL")

// ExifTool is the path to exiftool
var ExifTool = getEnv("EXIFTOOL")

// TargetICCProfileIIIF is the path to the target ICC profile
// for generation of pyramidal TIFFs for use by IIIF image server
var TargetICCProfileIIIF = getEnv("TARGET_ICC_PROFILE_IIIF")

// TargetICCProfileTIFF is the path to the target ICC profile
// for generation of downloadable TIFFs
var TargetICCProfileTIFF = getEnv("TARGET_ICC_PROFILE_TIFF")

/**************************
 * END Command-line paths *
 **************************/

// Returns the value of the environment variable named name
func getEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Printf("WARNING env var %s is not set\n", name)
	}
	return value
}
