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

// SanitizeDir allows a directory source to be set from multiple values. If any
// value is defined, that value is used. If none are defined, the current
// working directory is retrieved.
func SanitizeDir(dirs ...string) (string, error) {
	for _, dir := range dirs {
		if len(dir) > 0 {
			return dir, nil
		}
	}

	cwd, err := getwd()
	if err != nil {
		return "", fmt.Errorf("cannot determine cwd: %v", err)
	}

	return cwd, nil
}
