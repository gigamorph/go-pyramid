package util

import (
	"os/exec"
)

// GetStderr tries to extract what was printed to stderr from exec,
// assuming e is an exec.ExitError.
func GetStderr(e error) string {
	e2, ok := e.(*exec.ExitError)
	if (!ok) {
		return ""
	}
	return string(e2.Stderr)
}
