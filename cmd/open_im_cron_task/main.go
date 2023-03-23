package main

import (
	"Open_IM/internal/cron_task"
	"flag"
	"fmt"
	"time"
)

func main() {
	var userID = flag.String("userID", "", "userID to clear msg and reset seq")
	var workingGroupID = flag.String("workingGroupID", "", "workingGroupID to clear msg and reset seq")
	flag.Parse()
	fmt.Println(time.Now(), "start cronTask", *userID, *workingGroupID)
	cronTask.StartCronTask(*userID, *workingGroupID)
}
