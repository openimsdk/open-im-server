package main

import (
	"OpenIM/internal/tools"
	"OpenIM/pkg/common/cmd"
	"fmt"
	"os"
)

func main() {
	cronTaskCmd := cmd.NewCronTaskCmd()
	cronTaskCmd.AddRunE(tools.StartCronTask)
	if err := cronTaskCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
