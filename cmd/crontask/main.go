package main

import (
	"OpenIM/internal/tools"
	"OpenIM/pkg/common/cmd"
	"fmt"
	"os"
)

func main() {
	cronTaskCmd := cmd.NewCronTaskCmd()
	if err := cronTaskCmd.Exec(tools.StartCronTask); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
