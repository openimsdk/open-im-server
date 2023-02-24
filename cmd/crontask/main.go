package main

import (
	"OpenIM/internal/task"
	"OpenIM/pkg/common/config"
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now(), "start cronTask")
	if err := config.InitConfig(); err != nil {
		panic(err.Error())
	}
	if err := task.StartCronTask(); err != nil {
		panic(err.Error())
	}
}
