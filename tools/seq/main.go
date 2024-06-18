package main

import (
	"flag"
	"fmt"
	"github.com/openimsdk/open-im-server/v3/tools/seq/internal"
)

func main() {
	var config string
	flag.StringVar(&config, "redis", "/Users/chao/Desktop/project/open-im-server/config", "config directory")
	flag.Parse()
	if err := internal.Main(config); err != nil {
		fmt.Println("seq task", err)
	}
}
