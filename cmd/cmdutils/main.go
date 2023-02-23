package main

import (
	"OpenIM/internal/task"
	"OpenIM/pkg/common/config"
	"flag"
	"fmt"
	"time"
)

func main() {
	var userID = flag.String("user_id", "", "userID to clear msg and reset seq")
	var superGroupID = flag.String("super_group_id", "", "superGroupID to clear msg and reset seq")
	var fixAllSeq = flag.Bool("fix_all_seq", false, "fix seq")
	var configPath = flag.String("config_path", "../config/", "config folder")
	flag.Parse()
	if err := config.InitConfig(*configPath); err != nil {
		panic(err.Error())
	}
	fmt.Println(time.Now(), "start cronTask", *userID, *superGroupID)
	task.FixSeq(*userID, *superGroupID, *fixAllSeq)
}
