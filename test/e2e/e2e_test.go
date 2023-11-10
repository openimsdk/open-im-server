package e2e

import (
	"flag"
	"testing"

	"github.com/openimsdk/open-im-server/v3/test/e2e/framework/config"
)

// handleFlags sets up all flags and parses the command line.
func handleFlags() {
	config.CopyFlags(config.Flags, flag.CommandLine)
	flag.Parse()
}

func TestMain(m *testing.M) {
	handleFlags()
	m.Run()
}

func TestE2E(t *testing.T) {
	RunE2ETests(t)
}