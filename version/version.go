package version

import (
	_ "embed"
	"strings"
)

//go:embed version
var Version string

func init() {
	Version = strings.Trim(Version, "\n")
	Version = strings.TrimSpace(Version)
}
