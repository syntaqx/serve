package config

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

var getwd = os.Getwd

// Config holds the runtime configuration for the serve binary. Every field can
// be set from a CLI flag or its corresponding SERVE_* environment variable,
// with flags taking precedence over the environment.
type Config struct {
	Debug            bool
	Host             string
	Port             string
	EnableSSL        bool
	CertFile         string
	KeyFile          string
	Directory        string
	UsersFile        string
	ShowHidden       bool
	DirectoryListing bool
	HealthPath       string
	CORSOrigin       string
	LogFormat        string
	LogLevel         string
	ShowVersion      bool
}

// Load parses the given arguments (typically os.Args[1:]) into a Config,
// resolving defaults from the environment. It returns the config along with any
// remaining positional arguments. Flag parsing output is written to out.
func Load(args []string, out io.Writer) (*Config, []string, error) {
	c := &Config{}

	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(out)

	fs.BoolVar(&c.Debug, "debug", envBool("SERVE_DEBUG", false), "enable debug logging")
	fs.StringVar(&c.Host, "host", envStr("SERVE_HOST", "0.0.0.0"), "host address to bind to")
	fs.StringVar(&c.Port, "port", envPort(), "listening port (also honors PORT)")
	fs.BoolVar(&c.EnableSSL, "ssl", envBool("SERVE_SSL", false), "enable https")
	fs.StringVar(&c.CertFile, "cert", envStr("SERVE_CERT", "cert.pem"), "path to the TLS certificate file")
	fs.StringVar(&c.KeyFile, "key", envStr("SERVE_KEY", "key.pem"), "path to the TLS key file")
	fs.StringVar(&c.Directory, "dir", envStr("SERVE_DIR", ""), "directory to serve (defaults to the first argument or the current directory)")
	fs.StringVar(&c.UsersFile, "users", envStr("SERVE_USERS", "users.dat"), "path to a BasicAuth users file (user:password or user:bcrypt-hash per line)")
	fs.BoolVar(&c.ShowHidden, "all", envBool("SERVE_ALL", false), "serve dotfiles, which are hidden by default")
	fs.BoolVar(&c.DirectoryListing, "dirlisting", envBool("SERVE_DIRLISTING", false), "enable an automatic listing for directories without an index.html (disabled by default)")
	fs.StringVar(&c.HealthPath, "health-path", envStr("SERVE_HEALTH_PATH", "/healthz"), "path for the health-check endpoint (empty to disable)")
	fs.StringVar(&c.CORSOrigin, "cors-origin", envStr("SERVE_CORS_ORIGIN", "*"), "value for the Access-Control-Allow-Origin header")
	fs.StringVar(&c.LogFormat, "log-format", envStr("SERVE_LOG_FORMAT", "text"), "log output format: text or json")
	fs.StringVar(&c.LogLevel, "log-level", envStr("SERVE_LOG_LEVEL", "info"), "log level: debug, info, warn, or error")
	fs.BoolVar(&c.ShowVersion, "version", false, "print version information and exit")

	if err := fs.Parse(args); err != nil {
		return nil, nil, err
	}

	return c, fs.Args(), nil
}

// SanitizeDir returns the first non-empty directory from the provided values,
// falling back to the current working directory when none are set.
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

func envStr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func envBool(key string, def bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

// envPort resolves the listening port, preferring SERVE_PORT, then the
// conventional PORT variable (used by Heroku and others), then 8080.
func envPort() string {
	if v, ok := os.LookupEnv("SERVE_PORT"); ok {
		return v
	}
	if v, ok := os.LookupEnv("PORT"); ok {
		return v
	}
	return "8080"
}
