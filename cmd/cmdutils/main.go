package main

import (
	"Open_IM/internal/task"
	"flag"
	"fmt"
	"time"
)

func main() {
	var userID = flag.String("userID", "", "userID to clear msg and reset seq")
	var workingGroupID = flag.String("workingGroupID", "", "workingGroupID to clear msg and reset seq")
	var fixAllSeq = flag.Bool("fixAllSeq", false, "fix seq")
	flag.Parse()
	fmt.Println(time.Now(), "start cronTask", *userID, *workingGroupID)
	task.FixSeq(*userID, *workingGroupID, *fixAllSeq)
}
