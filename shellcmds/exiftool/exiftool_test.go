package exiftool

import (
	"testing"

	"github.com/gigamorph/go-pyramid/util"
	"github.com/stretchr/testify/assert"
)

func TestAddTags(t *testing.T) {
	testFile := "../../tmp/a.tif"
	_, err := util.CopyFile("../../test/resources/images/ba-obj-5005-0004-pub-print-lg.tif", testFile)
	if err != nil {
		panic(err)
	}
	caption := "I will arise and go now"
	copyrightNotice := "Anybody's guess"
	imageCredit := "Institute of Pagragon of Aethetics"
	webRightsStatement := "http://example.org/statement"
	usageTerms := "http://example.org/usage"
	copyrightStatus := "False"
	source := "Museum of Compassion"

	AddTags(testFile, TagsInput{
		CopyrightNotice:    copyrightNotice,
		ImageCredit:        imageCredit,
		WebRightsStatement: webRightsStatement,
		UsageTerms:         usageTerms,
		Caption:            caption,
		CopyrightStatus:    copyrightStatus,
		Source:             source,
	})

	tagValue, err := GetTag(testFile, "MWG:copyright")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, copyrightNotice, tagValue, "Add copyright notice")

	tagValue, err = GetTag(testFile, "XMP-photoshop:Credit")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, imageCredit, tagValue, "Add image credit (tag 1)")

	tagValue, err = GetTag(testFile, "credit")
	if err != nil {
		panic(err)
	}
	// IPTC Credit is limited to 32 characters
	assert.Equal(t, imageCredit[:32], tagValue, "Add image credit (tag 2)")

	tagValue, err = GetTag(testFile, "xmp:webstatement")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, webRightsStatement, tagValue, "Add web rights statement (tag 1)")

	tagValue, err = GetTag(testFile, "photoshop:URL")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, webRightsStatement, tagValue, "Add web rights statement (tag 2)")

	tagValue, err = GetTag(testFile, "usageterms")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, usageTerms, tagValue, "Add usage terms")

	tagValue, err = GetTag(testFile, "MWG:description")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, caption, tagValue, "Add caption")

	tagValue, err = GetTag(testFile, "XMP-xmpRights:marked")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, copyrightStatus, tagValue, "Add copyright status")

	tagValue, err = GetTag(testFile, "XMP-photoshop:Source")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, source, tagValue, "Add source (tag 1)")

	tagValue, err = GetTag(testFile, "iptc:source")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, source, tagValue, "Add source (tag 2)")
}
