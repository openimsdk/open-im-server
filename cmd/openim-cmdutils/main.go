package main

import (
	"github.com/openimsdk/open-im-server/v3/pkg/common/cmd"
	"github.com/openimsdk/tools/system/program"
)

func main() {
	msgUtilsCmd := cmd.NewMsgUtilsCmd("openIMCmdUtils", "openIM cmd utils", nil)
	getCmd := cmd.NewGetCmd()
	fixCmd := cmd.NewFixCmd()
	clearCmd := cmd.NewClearCmd()
	seqCmd := cmd.NewSeqCmd()
	msgCmd := cmd.NewMsgCmd()
	getCmd.AddCommand(seqCmd.GetSeqCmd(), msgCmd.GetMsgCmd())
	getCmd.AddSuperGroupIDFlag()
	getCmd.AddUserIDFlag()
	getCmd.AddConfigDirFlag()
	getCmd.AddIndexFlag()
	getCmd.AddBeginSeqFlag()
	getCmd.AddLimitFlag()
	// openIM get seq --userID=xxx
	// openIM get seq --superGroupID=xxx
	// openIM get msg --userID=xxx --beginSeq=100 --limit=10
	// openIM get msg --superGroupID=xxx --beginSeq=100 --limit=10

	fixCmd.AddCommand(seqCmd.FixSeqCmd())
	fixCmd.AddSuperGroupIDFlag()
	fixCmd.AddUserIDFlag()
	fixCmd.AddConfigDirFlag()
	fixCmd.AddIndexFlag()
	fixCmd.AddFixAllFlag()
	// openIM fix seq --userID=xxx
	// openIM fix seq --superGroupID=xxx
	// openIM fix seq --fixAll

	clearCmd.AddCommand(msgCmd.ClearMsgCmd())
	clearCmd.AddSuperGroupIDFlag()
	clearCmd.AddUserIDFlag()
	clearCmd.AddConfigDirFlag()
	clearCmd.AddIndexFlag()
	clearCmd.AddClearAllFlag()
	clearCmd.AddBeginSeqFlag()
	clearCmd.AddLimitFlag()
	// openIM clear msg --userID=xxx --beginSeq=100 --limit=10
	// openIM clear msg --superGroupID=xxx --beginSeq=100 --limit=10
	// openIM clear msg --clearAll
	msgUtilsCmd.AddCommand(&getCmd.Command, &fixCmd.Command, &clearCmd.Command)
	if err := msgUtilsCmd.Execute(); err != nil {
		program.ExitWithError(err)
	}
}
