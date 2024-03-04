package main

import (
	"log"

	"github.com/openimsdk/open-im-server/tools/codescan/checker"
	"github.com/openimsdk/open-im-server/tools/codescan/config"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	err = checker.WalkDirAndCheckComments(cfg)
	if err != nil {
		panic(err)
	}
}
