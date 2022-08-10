package main

import (
	"Open_IM/internal/cron_task"
	"fmt"
)

func main() {
	fmt.Println("start cronTask")
	cronTask.StartCronTask()
}
