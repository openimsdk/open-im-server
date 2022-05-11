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

func GetUploadAppNewName(appType int, version, fileName, yamlName string) (string, string, error) {
	var newFileName, newYamlName = "_" + version + "_app", "_" + version + "_yaml"
	switch appType {
	case constant.IOSPlatformID:
		newFileName = constant.IOSPlatformStr + newFileName
		newYamlName = constant.IOSPlatformStr + newYamlName
	case constant.AndroidPlatformID:
		newFileName = constant.AndroidPlatformStr + newFileName
		newYamlName = constant.AndroidPlatformStr + newYamlName
	case constant.WindowsPlatformID:
		newFileName = constant.WindowsPlatformStr + newFileName
		newYamlName = constant.WindowsPlatformStr + newYamlName
	case constant.OSXPlatformID:
		newFileName = constant.OSXPlatformStr + newFileName
		newYamlName = constant.OSXPlatformStr + newYamlName
	case constant.WebPlatformID:
		newFileName = constant.WebPlatformStr + newFileName
		newYamlName = constant.WebPlatformStr + newYamlName
	case constant.MiniWebPlatformID:
		newFileName = constant.MiniWebPlatformStr + newFileName
		newYamlName = constant.MiniWebPlatformStr + newYamlName
	case constant.LinuxPlatformID:
		newFileName = constant.LinuxPlatformStr + newFileName
		newYamlName = constant.LinuxPlatformStr + newYamlName
	default:
		return "", "", errors.New("invalid app type")
	}
	suffixFile := path.Ext(fileName)
	suffixYaml := path.Ext(yamlName)
	newFileName = fmt.Sprintf("%s%s", newFileName, suffixFile)
	newYamlName = fmt.Sprintf("%s%s", newYamlName, suffixYaml)
	if yamlName == "" {
		newYamlName = ""
	}
	return newFileName, newYamlName, nil
}
