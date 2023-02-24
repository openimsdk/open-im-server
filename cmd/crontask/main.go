package main

import (
	"OpenIM/internal/task"
	"OpenIM/pkg/common/config"
	"flag"
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now(), "start cronTask")
	var configPath = flag.String("config_path", "../config/", "config folder")
	flag.Parse()
	if err := config.InitConfig(*configPath); err != nil {
		panic(err.Error())
	}
	if err := task.StartCronTask(); err != nil {
		panic(err.Error())
	}
}
