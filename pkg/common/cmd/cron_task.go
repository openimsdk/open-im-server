package cmd

import "github.com/spf13/cobra"

type CronTaskCmd struct {
	*RootCmd
}

func NewCronTaskCmd() *CronTaskCmd {
	return &CronTaskCmd{NewRootCmd()}
}

func (c *CronTaskCmd) AddRunE(f func() error) {
	c.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return f()
	}
}
