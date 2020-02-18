package util

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// Exec is a utility wrapper around exec.Command.
func Exec(command string, args []string) (string, error) {
	log.Printf("util.Exec %s %s", command, strings.Join(args, " "))
	out, err := exec.Command(command, args...).Output()
	if err != nil {
		return string(out), fmt.Errorf("util.Exec exec.Command failed - %v - %s", err, GetStderr(err))
	}
	sout := strings.TrimSpace(string(out))
	return sout, nil
}
