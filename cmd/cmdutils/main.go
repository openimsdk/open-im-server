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
	// clear msg by id
	var userIDClearMsg = flag.String("user_id_fix_seq", "", "userID to clear msg and reset seq")
	var superGroupIDClearMsg = flag.String("super_group_id_fix_seq", "", "superGroupID to clear msg and reset seq")
	// show seq by id
	var userIDShowSeq = flag.String("user_id_show_seq", "", "show userID")
	var superGroupIDShowSeq = flag.String("super_group_id_show_seq", "", "userID to clear msg and reset seq")
	// fix seq by id
	var userIDFixSeq = flag.String("user_id_fix_seq", "", "userID to Fix Seq")
	var superGroupIDFixSeq = flag.String("super_group_id_fix_seq", "", "super groupID to fix Seq")
	var fixAllSeq = flag.Bool("fix_all_seq", false, "fix seq")

	flag.Parse()
	fmt.Println(time.Now(), "start cronTask", *userIDFixSeq, *superGroupIDFixSeq)
	task.FixSeq(*userID, *superGroupID, *fixAllSeq)
}
