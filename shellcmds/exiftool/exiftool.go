package exiftool

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gigamorph/go-pyramid/config"
	"github.com/gigamorph/go-pyramid/util"
)

// GetTag extracts a tag value from the image file
func GetTag(filePath, tagName string) (string, error) {
	var out string

	args := []string{
		"-TAG", fmt.Sprintf("-%s", tagName),
		filePath,
	}

	out, err := util.Exec(config.ExifTool, args)
	if err != nil {
		return "", fmt.Errorf("exiftool.GetTag failed - %v", err)
	}

	r := regexp.MustCompile(`^[^:]+: (.*)$`)
	m := r.FindStringSubmatch(out)
	if len(m) < 1 {
		return "", nil
	}
	return strings.TrimSpace(m[1]), nil
}

// AddTags invokes exiftool with the specified options to apply tags to the image file
func AddTags(filePath string, options TagsInput) (string, error) {
	var out string
	args := make([]string, 0, 8)

	if options.CopyrightNotice != "" {
		args = append(args, fmt.Sprintf("-MWG:copyright=%s", options.CopyrightNotice))
	}
	if options.ImageCredit != "" {
		args = append(args, fmt.Sprintf("-XMP-photoshop:Credit=%s", options.ImageCredit))
		args = append(args, fmt.Sprintf("-credit=%s", options.ImageCredit))
	}
	if options.WebRightsStatement != "" {
		args = append(args, fmt.Sprintf("-xmp:webstatement=%s", options.WebRightsStatement))
		args = append(args, fmt.Sprintf("-photoshop:URL=%s", options.WebRightsStatement))
	}
	if options.UsageTerms != "" {
		args = append(args, fmt.Sprintf("-usageterms=%s", options.UsageTerms))
	}
	if options.Caption != "" {
		args = append(args, fmt.Sprintf("-MWG:description=%s", options.Caption))
	}
	if options.CopyrightStatus != "" {
		args = append(args, fmt.Sprintf("-XMP-xmpRights:marked=%s", options.CopyrightStatus))
	}
	if options.Source != "" {
		args = append(args, fmt.Sprintf("-XMP-photoshop:Source=%s", options.Source))
		args = append(args, fmt.Sprintf("-iptc:source=%s", options.Source))
	}

	args = append(args, filePath)

	out, err := util.Exec(config.ExifTool, args)
	if err != nil {
		return "", fmt.Errorf("exiftool.Run failed - %v", err)
	}
	return out, nil
}
