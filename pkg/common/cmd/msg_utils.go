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
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/spf13/cobra"
)

type MsgUtilsCmd struct {
	cobra.Command
}

func (m *MsgUtilsCmd) AddUserIDFlag() {
	m.Command.PersistentFlags().StringP("userID", "u", "", "openIM userID")
}
func (m *MsgUtilsCmd) AddIndexFlag() {
	m.Command.PersistentFlags().IntP(config.FlagTransferIndex, "i", 0, "process startup sequence number")
}

func (m *MsgUtilsCmd) AddConfigDirFlag() {
	m.Command.PersistentFlags().StringP(config.FlagConf, "c", "", "path of config directory")

}

func (m *MsgUtilsCmd) getUserIDFlag(cmdLines *cobra.Command) string {
	userID, _ := cmdLines.Flags().GetString("userID")
	return userID
}

func (m *MsgUtilsCmd) AddFixAllFlag() {
	m.Command.PersistentFlags().BoolP("fixAll", "f", false, "openIM fix all seqs")
}

/* func (m *MsgUtilsCmd) getFixAllFlag(cmdLines *cobra.Command) bool {
	fixAll, _ := cmdLines.Flags().GetBool("fixAll")
	return fixAll
} */

func (m *MsgUtilsCmd) AddClearAllFlag() {
	m.Command.PersistentFlags().BoolP("clearAll", "", false, "openIM clear all seqs")
}

/* func (m *MsgUtilsCmd) getClearAllFlag(cmdLines *cobra.Command) bool {
	clearAll, _ := cmdLines.Flags().GetBool("clearAll")
	return clearAll
} */

func (m *MsgUtilsCmd) AddSuperGroupIDFlag() {
	m.Command.PersistentFlags().StringP("superGroupID", "g", "", "openIM superGroupID")
}

func (m *MsgUtilsCmd) getSuperGroupIDFlag(cmdLines *cobra.Command) string {
	superGroupID, _ := cmdLines.Flags().GetString("superGroupID")
	return superGroupID
}

func (m *MsgUtilsCmd) AddBeginSeqFlag() {
	m.Command.PersistentFlags().Int64P("beginSeq", "b", 0, "openIM beginSeq")
}

/* func (m *MsgUtilsCmd) getBeginSeqFlag(cmdLines *cobra.Command) int64 {
	beginSeq, _ := cmdLines.Flags().GetInt64("beginSeq")
	return beginSeq
} */

func (m *MsgUtilsCmd) AddLimitFlag() {
	m.Command.PersistentFlags().Int64P("limit", "l", 0, "openIM limit")
}

/* func (m *MsgUtilsCmd) getLimitFlag(cmdLines *cobra.Command) int64 {
	limit, _ := cmdLines.Flags().GetInt64("limit")
	return limit
} */

func (m *MsgUtilsCmd) Execute() error {
	return m.Command.Execute()
}

func NewMsgUtilsCmd(use, short string, args cobra.PositionalArgs) *MsgUtilsCmd {
	return &MsgUtilsCmd{
		Command: cobra.Command{
			Use:   use,
			Short: short,
			Args:  args,
		},
	}
}

type GetCmd struct {
	*MsgUtilsCmd
}

func NewGetCmd() *GetCmd {
	return &GetCmd{
		NewMsgUtilsCmd("get [resource]", "get action", cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs)),
	}
}

type FixCmd struct {
	*MsgUtilsCmd
}

func NewFixCmd() *FixCmd {
	return &FixCmd{
		NewMsgUtilsCmd("fix [resource]", "fix action", cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs)),
	}
}

type ClearCmd struct {
	*MsgUtilsCmd
}

func NewClearCmd() *ClearCmd {
	return &ClearCmd{
		NewMsgUtilsCmd("clear [resource]", "clear action", cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs)),
	}
}

type SeqCmd struct {
	*MsgUtilsCmd
}

func NewSeqCmd() *SeqCmd {
	seqCmd := &SeqCmd{
		NewMsgUtilsCmd("seq", "seq", nil),
	}
	return seqCmd
}

func (s *SeqCmd) GetSeqCmd() *cobra.Command {
	s.Command.Run = func(cmdLines *cobra.Command, args []string) {

	}
	return &s.Command
}

func (s *SeqCmd) FixSeqCmd() *cobra.Command {
	return &s.Command
}

type MsgCmd struct {
	*MsgUtilsCmd
}

func NewMsgCmd() *MsgCmd {
	msgCmd := &MsgCmd{
		NewMsgUtilsCmd("msg", "msg", nil),
	}
	return msgCmd
}

func (m *MsgCmd) GetMsgCmd() *cobra.Command {
	return &m.Command
}

func (m *MsgCmd) ClearMsgCmd() *cobra.Command {
	return &m.Command
}
