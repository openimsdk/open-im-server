package main

import (
	"flag"
	"fmt"
	"github.com/openimsdk/open-im-server/v3/tools/s3/internal"
	"os"
)

func main() {
	var (
		name   string
		config string
	)
	flag.StringVar(&name, "name", "", "old previous storage name")
	flag.StringVar(&config, "config", "", "config directory")
	flag.Parse()
	if err := internal.Main(config, name); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, "success")
}
