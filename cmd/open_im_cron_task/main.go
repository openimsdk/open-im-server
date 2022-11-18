package main

import (
	"Open_IM/internal/cron_task"
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now(), "start cronTask")
	cronTask.StartCronTask()
}
