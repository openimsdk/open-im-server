package main

import (
	"OpenIM/internal/tools"
	"OpenIM/pkg/common/config"
	"context"
	"flag"
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
	msgTool, err := tools.InitMsgTool()
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()
	if userIDFixSeq != nil {
		msgTool.GetAndFixUserSeqs(ctx, *userIDFixSeq)
	}
	if superGroupIDFixSeq != nil {
		msgTool.FixGroupSeq(ctx, *superGroupIDFixSeq)
	}
	if fixAllSeq != nil {
		msgTool.FixAllSeq(ctx)
	}
	if userIDClearMsg != nil {
		msgTool.ClearUsersMsg(ctx, []string{*userIDClearMsg})
	}

	if superGroupIDClearMsg != nil {
		msgTool.ClearSuperGroupMsg(ctx, []string{*superGroupIDClearMsg})
	}
	if userIDShowSeq != nil {
		msgTool.ShowUserSeqs(ctx, *userIDShowSeq)
	}

	if superGroupIDShowSeq != nil {
		msgTool.ShowSuperGroupSeqs(ctx, *superGroupIDShowSeq)
	}
}
