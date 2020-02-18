package agent

import (
	"fmt"
	"log"
	"os"

	"github.com/gigamorph/go-pyramid/config"
)

func init() {
	// Create temp directory (workspace)
	tempDir := config.TempDir
	log.Printf("Making sure the temporary workspace exists at %s", tempDir)
	err := os.MkdirAll(tempDir, 0700)
	if err != nil {
		panic(fmt.Errorf("Failed to create directory %s - %v", tempDir, err))
	}
}
