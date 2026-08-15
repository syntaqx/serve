package commands

import (
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
)

// Version implements the command `version` which outputs the current binary
// release version, if any, along with the platform, Go toolchain, and (when
// built from a VCS checkout) the source revision.
func Version(version string, w io.Writer) error {
	line := fmt.Sprintf("serve version %s %s/%s %s", version, runtime.GOOS, runtime.GOARCH, runtime.Version())

	if rev := vcsRevision(); rev != "" {
		line += fmt.Sprintf(" (%s)", rev)
	}

	_, err := fmt.Fprintln(w, line)
	return err
}

// readBuildInfo is a seam for tests to supply synthetic build information.
var readBuildInfo = debug.ReadBuildInfo

// vcsRevision returns a short commit hash embedded by the Go toolchain at build
// time, or an empty string when unavailable.
func vcsRevision() string {
	info, ok := readBuildInfo()
	if !ok {
		return ""
	}
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			if len(s.Value) > 12 {
				return s.Value[:12]
			}
			return s.Value
		}
	}
	return ""
}
