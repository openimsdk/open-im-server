package utils

import "os"

// Determine whether the given path is a folder
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// Determine whether the given path is a file
func IsFile(path string) bool {
	return !IsDir(path)
}

// Create a directory
func MkDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
