package main

import (
	"OpenIM/internal/tools"
	"OpenIM/pkg/common/cmd"
)

func main() {
	cronTaskCmd := cmd.NewCronTaskCmd()
	if err := cronTaskCmd.Exec(tools.StartCronTask); err != nil {
		panic(err.Error())
	}
}
