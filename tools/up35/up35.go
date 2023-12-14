package main

import (
	"flag"
	"github.com/openimsdk/open-im-server/v3/tools/up35/pkg"
	"log"
	"os"
)

func main() {
	var path string
	flag.StringVar(&path, "c", "", "path config file")
	flag.Parse()
	if err := pkg.Main(path); err != nil {
		log.Fatal(err)
		return
	}
	os.Exit(0)
}
