//go:build mage
// +build mage

package main

import (
	"os"

	"github.com/mo3et/openim-gomake/mageutil"
)

var Default = Build

func Build() {
	mageutil.Build()
}

func Start() {
	mageutil.InitForSSC()
	err := setMaxOpenFiles()
	if err != nil {
		mageutil.PrintRed("setMaxOpenFiles failed " + err.Error())
		os.Exit(1)
	}
	mageutil.StartToolsAndServices()
}

func Stop() {
	mageutil.StopAndCheckBinaries()
}

func Check() {
	mageutil.CheckAndReportBinariesStatus()
}
