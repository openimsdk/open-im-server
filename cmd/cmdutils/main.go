package main

import (
	"OpenIM/internal/task"
	"OpenIM/pkg/common/config"
	"flag"
	"fmt"
	"time"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err.Error())
	}
	var userID = flag.String("user_id", "", "userID to clear msg and reset seq")
	var superGroupID = flag.String("super_group_id", "", "superGroupID to clear msg and reset seq")
	var fixAllSeq = flag.Bool("fix_all_seq", false, "fix seq")
	flag.Parse()
	fmt.Println(time.Now(), "start cronTask", *userID, *superGroupID)
	task.FixSeq(*userID, *superGroupID, *fixAllSeq)
}
