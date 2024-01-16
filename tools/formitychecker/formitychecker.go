package main

import (
	"flag"
	"fmt"

	"github.com/openimsdk/open-im-server/tools/formitychecker/checker"
	"github.com/openimsdk/open-im-server/tools/formitychecker/config"
)

func main() {
	defaultTargetDirs := "."
	defaultIgnoreDirs := "components,.git"

	var targetDirs string
	var ignoreDirs string
	flag.StringVar(&targetDirs, "target", defaultTargetDirs, "Directories to check (default: current directory)")
	flag.StringVar(&ignoreDirs, "ignore", defaultIgnoreDirs, "Directories to ignore (default: A/, B/)")
	flag.Parse()

	conf := config.NewConfig(targetDirs, ignoreDirs)

	err := checker.CheckDirectory(conf)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
