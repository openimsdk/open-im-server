/**
为确保运行结果正常输出，请不要移动此文件
 */
package base

import (
	"errors"
	"path/filepath"
	"runtime"
)

// 获取项目基准目录
func baseDir(p string) (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("runtime.Caller error.")
	}
	dir:= filepath.Dir( filepath.Dir( filepath.Dir( file ) ) )
	return dir+ p, nil
}
// 获取项目根目录
func ProjectDir(p string) string {
	if p=="" {
		p= "/"
	}
	dir, err:= baseDir(p)
	if err!= nil {
		panic(err)
	}
	return dir
}
// 获取各个常用目录
func ConfigDir() string {
	return ProjectDir("/config/")
}
func ScriptDir() string {
	return ProjectDir("/script/")
}
func BinDir() string {
	return ProjectDir("/bin/")
}
func LogDir() string {
	return ProjectDir("/logs/")
}
func PkgDir() string {
	return ProjectDir("/pkg/")
}
func InternalDir() string {
	return ProjectDir("/internal/")
}

