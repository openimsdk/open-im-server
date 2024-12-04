//go:build mage
// +build mage

package main

import (
	"flag"
	"os"

	"github.com/openimsdk/gomake/mageutil"
)

var Default = Build

func Build() {
	flag.Parse()

	bin := flag.Args()
	if len(bin) != 0 {
		bin = bin[1:]
	}

	mageutil.Build(bin)
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
