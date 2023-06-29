package cmd

import "github.com/spf13/cobra"

type ApiCmd struct {
	*RootCmd
}

func NewApiCmd() *ApiCmd {
	return &ApiCmd{NewRootCmd("api")}
}

func (a *ApiCmd) AddApi(f func(port int) error) {
	a.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return f(a.getPortFlag(cmd))
	}
}
