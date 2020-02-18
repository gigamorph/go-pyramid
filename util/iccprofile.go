package util

import (
	"bytes"
	"encoding/binary"
	"log"
)

// TagSigDesc is the hexadecimal representation of the tag signature for "desc"
const TagSigDesc uint32 = 0x64657363

// GetICCProfileDesc extracts the ASCII name part ("desc") from the byte array of the ICC profile.
func GetICCProfileDesc(iccProfile []byte) string {
	var err error
	var nTags, tagSignature, tagDataOffset, tagDataSize uint32
	var start int

	err = binary.Read(bytes.NewReader(iccProfile[128:132]), binary.BigEndian, &nTags)
	if err != nil {
		log.Printf("Error reading tags count - %v", err)
		return ""
	}

	for i := 0; i < int(nTags); i++ {
		start = 132 + 12*i
		err = binary.Read(bytes.NewReader(iccProfile[start:start+4]), binary.BigEndian, &tagSignature)
		if err != nil {
			log.Printf("Error reading tag signature for %d - %v\n", i, err)
			break
		}
		if tagSignature == TagSigDesc {
			err = binary.Read(bytes.NewReader(iccProfile[start+4:start+8]), binary.BigEndian, &tagDataOffset)
			if err != nil {
				log.Printf("Error reading tag data offset for %d - %v\n", i, err)
				break
			}
			err = binary.Read(bytes.NewReader(iccProfile[start+8:start+12]), binary.BigEndian, &tagDataSize)
			if err != nil {
				log.Printf("Error reading tag data size for %d - %v\n", i, err)
				break
			}
			tagData := string(iccProfile[tagDataOffset : tagDataOffset+tagDataSize])
			if err != nil {
				log.Printf("Error reading tag data size for %d - %v\n", i, err)
				break

			}
			return tagData[4:]
		}
	}

	return ""
}
