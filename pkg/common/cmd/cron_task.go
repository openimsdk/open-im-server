package cmd

import "github.com/spf13/cobra"

type CronTaskCmd struct {
	*RootCmd
}

func NewCronTaskCmd() *CronTaskCmd {
	return &CronTaskCmd{NewRootCmd("cronTask")}
}

func (c *CronTaskCmd) addRunE(f func() error) {
	c.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return f()
	}
}

func (c *CronTaskCmd) Exec(f func() error) error {
	c.addRunE(f)
	return c.Execute()
}
