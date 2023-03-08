package main

import (
	"OpenIM/pkg/common/cmd"
	"fmt"
	"os"
)

func main() {
	msgUtilsCmd := cmd.NewMsgUtilsCmd()
	msgUtilsCmd.AddSuperGroupIDFlag()
	msgUtilsCmd.AddUserIDFlag()
	seqCmd := cmd.NewSeqCmd()
	msgCmd := cmd.NewMsgCmd()
	cmd.GetCmd.AddCommand(seqCmd.Command, msgCmd.Command)
	cmd.FixCmd.AddCommand(seqCmd.Command)
	cmd.GetCmd.AddCommand(msgCmd.Command)
	msgUtilsCmd.AddCommand(cmd.GetCmd, cmd.FixCmd, cmd.ClearCmd)
	if err := msgUtilsCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
