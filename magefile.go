//go:build mage
// +build mage

package main

import (
	"github.com/openimsdk/tools/utils/mageutil"
	"os"
	"strings"
)

var Default = Build

func Build() {
	platforms := os.Getenv("PLATFORMS")
	if platforms == "" {
		platforms = mageutil.DetectPlatform()
	}

	for _, platform := range strings.Split(platforms, " ") {
		mageutil.CompileForPlatform(platform)
	}

	mageutil.PrintGreen("Compilation complete.")
}

func Start() {
	setMaxOpenFiles()
	mageutil.StartToolsAndServices()
}

func Stop() {
	mageutil.StopAndCheckBinaries()
}

func Check() {
	mageutil.CheckAndReportBinariesStatus()
}
