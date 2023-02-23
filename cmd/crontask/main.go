package main

import (
	"Open_IM/internal/task"
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now(), "start cronTask")
	if err := task.StartCronTask(); err != nil {
		panic(err.Error())
	}
}
