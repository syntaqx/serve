package commands

import (
	"fmt"
	"io"
	"runtime"
)

// Version implements the command `version` which outputs the current binary
// release version, if any.
func Version(version string, w io.Writer) error {
	fmt.Fprintf(w, "serve version %s %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
	return nil
}
