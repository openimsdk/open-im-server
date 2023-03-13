package cmd

import (
	"OpenIM/internal/tools"
	"context"
	"github.com/spf13/cobra"
)

type MsgUtilsCmd struct {
	*RootCmd
	userID string

	superGroupID string

	clearAll bool

	fixAll bool
}

func NewMsgUtilsCmd() MsgUtilsCmd {
	return MsgUtilsCmd{RootCmd: NewRootCmd("msgUtils")}
}

func (m *MsgUtilsCmd) AddUserIDFlag() {
	m.Command.PersistentFlags().StringP("userID", "u", "", "openIM userID")
}

func (m *MsgUtilsCmd) GetUserIDFlag() string {
	return m.userID
}

func (m *MsgUtilsCmd) AddFixAllFlag() {
	m.Command.PersistentFlags().BoolP("fixAll", "c", false, "openIM fix all seqs")
}

func (m *MsgUtilsCmd) GetFixAllFlag() bool {
	return m.fixAll
}

func (m *MsgUtilsCmd) AddSuperGroupIDFlag() {
	m.Command.PersistentFlags().StringP("super-groupID", "u", "", "openIM superGroupID")
}

func (m *MsgUtilsCmd) GetSuperGroupIDFlag() string {
	return m.superGroupID
}

func (m *MsgUtilsCmd) AddClearAllFlag() bool {
	return m.clearAll
}

func (m *MsgUtilsCmd) GetClearAllFlag() bool {
	return m.clearAll
}

type SeqCmd struct {
	Command *cobra.Command
}

func (SeqCmd) RunCommand(cmdLines *cobra.Command, args []string) error {
	msgTool, err := tools.InitMsgTool()
	if err != nil {
		return err
	}
	userID, _ := cmdLines.Flags().GetString("userID")
	superGroupID, _ := cmdLines.Flags().GetString("superGroupID")
	fixAll, _ := cmdLines.Flags().GetBool("fixAll")
	ctx := context.Background()
	switch {
	case cmdLines.Parent() == GetCmd:
		switch {
		case userID != "":
			msgTool.ShowUserSeqs(ctx, userID)
		case superGroupID != "":
			msgTool.ShowSuperGroupSeqs(ctx, superGroupID)
		}
	case cmdLines.Parent() == FixCmd:
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
}

func NewSeqCmd() SeqCmd {
	seqCmd := SeqCmd{&cobra.Command{
		Use:   "seq",
		Short: "seq operation",
	}}
	seqCmd.Command.Flags().BoolP("fixAll", "c", false, "openIM fix all seqs")
	seqCmd.Command.RunE = seqCmd.RunCommand
	return seqCmd
}

type MsgCmd struct {
	Command *cobra.Command
}

func NewMsgCmd() MsgCmd {
	msgCmd := MsgCmd{&cobra.Command{
		Use:   "msg",
		Short: "msg operation",
	}}
	msgCmd.Command.RunE = msgCmd.RunCommand
	msgCmd.Command.Flags().BoolP("clearAll", "c", false, "openIM clear all timeout msgs")
	return msgCmd
}

func (*MsgCmd) RunCommand(cmdLines *cobra.Command, args []string) error {
	msgTool, err := tools.InitMsgTool()
	if err != nil {
		return err
	}
	userID, _ := cmdLines.Flags().GetString("userID")
	superGroupID, _ := cmdLines.Flags().GetString("superGroupID")
	clearAll, _ := cmdLines.Flags().GetBool("clearAll")
	ctx := context.Background()
	switch {
	case cmdLines.Parent() == GetCmd:
		switch {
		case userID != "":
			msgTool.ShowUserSeqs(ctx, userID)
		case superGroupID != "":
			msgTool.ShowSuperGroupSeqs(ctx, superGroupID)
		}
	case cmdLines.Parent() == ClearCmd:
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
}

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "get operation",
}

var FixCmd = &cobra.Command{
	Use:   "fix",
	Short: "fix seq operation",
}

var ClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "clear operation",
}
