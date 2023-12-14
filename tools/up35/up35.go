package main

import (
	"flag"
	"log"
	"os"

	"github.com/openimsdk/open-im-server/v3/tools/up35/pkg"
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
