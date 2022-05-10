package utils

import (
	"Open_IM/pkg/common/constant"
	"errors"
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

func GetUploadAppNewName(appType int, version string) (string, string, error) {
	var newFileName, newYamlName = version + "_app_", version + "_yaml_"
	switch appType {
	case constant.IOSPlatformID:
		newFileName += constant.IOSPlatformStr
		newYamlName += constant.IOSPlatformStr
	case constant.AndroidPlatformID:
		newFileName += constant.AndroidPlatformStr
		newYamlName += constant.AndroidPlatformStr
	case constant.WindowsPlatformID:
		newFileName += constant.WindowsPlatformStr
		newYamlName += constant.WindowsPlatformStr
	case constant.OSXPlatformID:
		newFileName += constant.OSXPlatformStr
		newYamlName += constant.OSXPlatformStr
	case constant.WebPlatformID:
		newFileName += constant.WebPlatformStr
		newYamlName += constant.WebPlatformStr
	case constant.MiniWebPlatformID:
		newFileName += constant.MiniWebPlatformStr
		newYamlName += constant.MiniWebPlatformStr
	case constant.LinuxPlatformID:
		newFileName += constant.LinuxPlatformStr
		newYamlName += constant.LinuxPlatformStr
	default:
		return "", "", errors.New("invalid app type")
	}
	return newFileName, newYamlName, nil
}
