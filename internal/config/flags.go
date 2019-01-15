package config

import (
	"fmt"
	"os"
)

var getwd = os.Getwd

// Flags are the expose configuration flags available to the serve binary.
type Flags struct {
	Host string
	Port int
	Dir  string
}

func SanitizeDir(optval, cmdval string) (string, error) {
	if len(optval) != 0 {
		return optval, nil
	} else if len(cmdval) != 0 {
		return cmdval, nil
	}

	cwd, err := getwd()
	if err != nil {
		return "", fmt.Errorf("cannot determine cwd: %v", err)
	}

	return cwd, nil
}
