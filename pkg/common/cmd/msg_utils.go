// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/OpenIMSDK/Open-IM-Server/internal/tools"
	"github.com/spf13/cobra"
)

// define a message util command struct
type MsgUtilsCmd struct {
	cobra.Command
	msgTool *tools.MsgTool
}

// add userID flag
func (m *MsgUtilsCmd) AddUserIDFlag() {
	m.Command.PersistentFlags().StringP("userID", "u", "", "openIM userID")
}

// get userID flag
func (m *MsgUtilsCmd) getUserIDFlag(cmdLines *cobra.Command) string {
	userID, _ := cmdLines.Flags().GetString("userID")
	return userID
}

// add fix all flag
func (m *MsgUtilsCmd) AddFixAllFlag() {
	m.Command.PersistentFlags().BoolP("fixAll", "f", false, "openIM fix all seqs")
}

// get fix all flag
func (m *MsgUtilsCmd) getFixAllFlag(cmdLines *cobra.Command) bool {
	fixAll, _ := cmdLines.Flags().GetBool("fixAll")
	return fixAll
}

// add clear all flag
func (m *MsgUtilsCmd) AddClearAllFlag() {
	m.Command.PersistentFlags().BoolP("clearAll", "c", false, "openIM clear all seqs")
}

// get clear all flag
func (m *MsgUtilsCmd) getClearAllFlag(cmdLines *cobra.Command) bool {
	clearAll, _ := cmdLines.Flags().GetBool("clearAll")
	return clearAll
}

// add super groupID flag
func (m *MsgUtilsCmd) AddSuperGroupIDFlag() {
	m.Command.PersistentFlags().StringP("superGroupID", "g", "", "openIM superGroupID")
}

// get super groupID flag
func (m *MsgUtilsCmd) getSuperGroupIDFlag(cmdLines *cobra.Command) string {
	superGroupID, _ := cmdLines.Flags().GetString("superGroupID")
	return superGroupID
}

// add begin sequence flag
func (m *MsgUtilsCmd) AddBeginSeqFlag() {
	m.Command.PersistentFlags().Int64P("beginSeq", "b", 0, "openIM beginSeq")
}

// get begin sequence flag
func (m *MsgUtilsCmd) getBeginSeqFlag(cmdLines *cobra.Command) int64 {
	beginSeq, _ := cmdLines.Flags().GetInt64("beginSeq")
	return beginSeq
}

// add limited flag
func (m *MsgUtilsCmd) AddLimitFlag() {
	m.Command.PersistentFlags().Int64P("limit", "l", 0, "openIM limit")
}

// get limited flag
func (m *MsgUtilsCmd) getLimitFlag(cmdLines *cobra.Command) int64 {
	limit, _ := cmdLines.Flags().GetInt64("limit")
	return limit
}

// execute
func (m *MsgUtilsCmd) Execute() error {
	return m.Command.Execute()
}

// new message utils command
func NewMsgUtilsCmd(use, short string, args cobra.PositionalArgs) *MsgUtilsCmd {
	return &MsgUtilsCmd{
		Command: cobra.Command{
			Use:   use,
			Short: short,
			Args:  args,
		},
	}
}

// define a getcommand dtruct
type GetCmd struct {
	*MsgUtilsCmd
}

// create a new command
func NewGetCmd() *GetCmd {
	return &GetCmd{
		NewMsgUtilsCmd("get [resource]", "get action", cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs)),
	}
}

// define a fix command struct
type FixCmd struct {
	*MsgUtilsCmd
}

// new a fixed command
func NewFixCmd() *FixCmd {
	return &FixCmd{
		NewMsgUtilsCmd("fix [resource]", "fix action", cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs)),
	}
}

// define a clear command
type ClearCmd struct {
	*MsgUtilsCmd
}

// create a new command
func NewClearCmd() *ClearCmd {
	return &ClearCmd{
		NewMsgUtilsCmd("clear [resource]", "clear action", cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs)),
	}
}

// define a sequnce command struct
type SeqCmd struct {
	*MsgUtilsCmd
}

// create a new seq command
func NewSeqCmd() *SeqCmd {
	seqCmd := &SeqCmd{
		NewMsgUtilsCmd("seq", "seq", nil),
	}
	return seqCmd
}

// get a sequence command
func (s *SeqCmd) GetSeqCmd() *cobra.Command {
	s.Command.Run = func(cmdLines *cobra.Command, args []string) {
		_, err := tools.InitMsgTool()
		if err != nil {
			panic(err)
		}
		userID := s.getUserIDFlag(cmdLines)
		superGroupID := s.getSuperGroupIDFlag(cmdLines)
		// beginSeq := s.getBeginSeqFlag(cmdLines)
		// limit := s.getLimitFlag(cmdLines)
		if userID != "" {
			// seq, err := msgTool.s(context.Background(), userID)
			if err != nil {
				panic(err)
			}
			// println(seq)
		} else if superGroupID != "" {
			// seq, err := msgTool.GetSuperGroupSeq(context.Background(), superGroupID)
			if err != nil {
				panic(err)
			}
			// println(seq)
		}
	}
	return &s.Command
}

// fix a sequence command
func (s *SeqCmd) FixSeqCmd() *cobra.Command {
	return &s.Command
}

// define a message command
type MsgCmd struct {
	*MsgUtilsCmd
}

// create a message command
func NewMsgCmd() *MsgCmd {
	msgCmd := &MsgCmd{
		NewMsgUtilsCmd("msg", "msg", nil),
	}
	return msgCmd
}

// get message command
func (m *MsgCmd) GetMsgCmd() *cobra.Command {
	return &m.Command
}

// clear message command
func (m *MsgCmd) ClearMsgCmd() *cobra.Command {
	return &m.Command
}
