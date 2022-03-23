package utils

import (
	"Open_IM/pkg/common/constant"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"
)

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

func GetNewFileNameAndContentType(fileName string, fileType int) (string, string) {
	suffix := path.Ext(fileName)
	newName := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Int(), fileName)
	contentType := ""
	if fileType == constant.ImageType {
		contentType = "image/" + suffix[1:]
	}
	return newName, contentType
}
