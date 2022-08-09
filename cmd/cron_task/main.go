package cron_task

import (
	"Open_IM/internal/cron_task"
	"fmt"
)

func main() {
	fmt.Println("start cronTask")
	cronTask.StartCronTask()
}
