package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/openimsdk/open-im-server/v3/tools/seq/internal"
)

func main() {
	var (
		config string
		second int
	)
	flag.StringVar(&config, "c", "", "config directory")
	flag.IntVar(&second, "sec", 3600*24, "delayed deletion of the original seq key after conversion")
	flag.Parse()
	if err := internal.Main(config, time.Duration(second)*time.Second); err != nil {
		fmt.Println("seq task", err)
		os.Exit(1)
		return
	}
	fmt.Println("seq task success!")
}
