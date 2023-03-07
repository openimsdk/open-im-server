package main

import (
	"OpenIM/internal/tools"
	"OpenIM/pkg/common/cmd"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var seqCmd = &cobra.Command{
	Use:   "seq",
	Short: "seq operation",
	RunE: func(cmdLines *cobra.Command, args []string) error {
		msgTool, err := tools.InitMsgTool()
		if err != nil {
			return err
		}
		userID, _ := cmdLines.Flags().GetString("userID")
		superGroupID, _ := cmdLines.Flags().GetString("superGroupID")
		fixAll, _ := cmdLines.Flags().GetBool("fixAll")
		ctx := context.Background()
		switch {
		case cmdLines.Parent() == getCmd:
			switch {
			case userID != "":
				msgTool.ShowUserSeqs(ctx, userID)
			case superGroupID != "":
				msgTool.ShowSuperGroupSeqs(ctx, superGroupID)
			}
		case cmdLines.Parent() == fixCmd:
			switch {
			case userID != "":
				_, _, err = msgTool.GetAndFixUserSeqs(ctx, userID)
			case superGroupID != "":
				err = msgTool.FixGroupSeq(ctx, userID)
			case fixAll:
				err = msgTool.FixAllSeq(ctx)
			}
		}
		return err
	},
}

var msgCmd = &cobra.Command{
	Use:   "msg",
	Short: "msg operation",
	RunE: func(cmdLines *cobra.Command, args []string) error {
		msgTool, err := tools.InitMsgTool()
		if err != nil {
			return err
		}
		userID, _ := cmdLines.Flags().GetString("userID")
		superGroupID, _ := cmdLines.Flags().GetString("superGroupID")
		clearAll, _ := cmdLines.Flags().GetBool("clearAll")
		ctx := context.Background()
		switch {
		case cmdLines.Parent() == getCmd:
			switch {
			case userID != "":
				msgTool.ShowUserSeqs(ctx, userID)
			case superGroupID != "":
				msgTool.ShowSuperGroupSeqs(ctx, superGroupID)
			}
		case cmdLines.Parent() == clearCmd:
			switch {
			case userID != "":
				msgTool.ClearUsersMsg(ctx, []string{userID})
			case superGroupID != "":
				msgTool.ClearSuperGroupMsg(ctx, []string{superGroupID})
			case clearAll:
				msgTool.AllUserClearMsgAndFixSeq()
			}
		}
		return nil
	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get operation",
}

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "fix seq operation",
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "clear operation",
}

func main() {
	cmd.RootCmd.PersistentFlags().StringP("userID", "u", "", "openIM userID")
	cmd.RootCmd.PersistentFlags().StringP("groupID", "u", "", "openIM superGroupID")
	seqCmd.Flags().BoolP("fixAll", "c", false, "openIM fix all seqs")
	msgCmd.Flags().BoolP("clearAll", "c", false, "openIM clear all timeout msgs")
	cmd.RootCmd.AddCommand(getCmd, fixCmd, clearCmd)
	getCmd.AddCommand(seqCmd, msgCmd)
	fixCmd.AddCommand(seqCmd)
	clearCmd.AddCommand(msgCmd)
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
