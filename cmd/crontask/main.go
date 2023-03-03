package main

import (
	"OpenIM/internal/task"
	"OpenIM/internal/tools"
	"OpenIM/pkg/common/config"
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now(), "start cronTask")
	if err := config.InitConfig(); err != nil {
		panic(err.Error())
	}
	if err := tools.StartCronTask(); err != nil {
		panic(err.Error())
	}
}
