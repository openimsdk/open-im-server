package cmd

type MsgUtilsCmd struct {
	RootCmd
	userID     string
	userIDFlag bool

	superGroupID     string
	superGroupIDFlag bool

	clearAll     bool
	clearAllFlag bool

	fixAll     bool
	fixAllFlag bool
}

func NewMsgUtilsCmd() MsgUtilsCmd {
	return MsgUtilsCmd{RootCmd: NewRootCmd()}
}

func (m *MsgUtilsCmd) AddUserIDFlag() {
	m.Command.PersistentFlags().StringP("userID", "u", "", "openIM userID")
	m.userIDFlag = true
}

func (m *MsgUtilsCmd) getUserIDFlag() {
	m.Command.PersistentFlags().StringP("userID", "u", "", "openIM userID")

}

func (m *MsgUtilsCmd) AddGroupIDFlag() {
	m.Command.PersistentFlags().StringP("super-groupID", "u", "", "openIM superGroupID")
}

func (m *MsgUtilsCmd) AddClearAllFlag() {
	m.Command.PersistentFlags().BoolP("clearAll", "c", false, "openIM clear all timeout msgs")
}

func (m *MsgUtilsCmd) AddFixAllFlag() {
	m.Command.PersistentFlags().BoolP("fixAll", "c", false, "openIM fix all seqs")
}
