package config

import (
	"strings"
)

// Config holds all the configuration parameters for the checker.
type Config struct {
	TargetDirs []string // Directories to check
	IgnoreDirs []string // Directories to ignore
}

// NewConfig creates and returns a new Config instance.
func NewConfig(targetDirs, ignoreDirs string) *Config {
	return &Config{
		TargetDirs: parseDirs(targetDirs),
		IgnoreDirs: parseDirs(ignoreDirs),
	}
}

// parseDirs splits a comma-separated string into a slice of directory names.
func parseDirs(dirs string) []string {
	if dirs == "" {
		return nil
	}
	return strings.Split(dirs, ",")
}
