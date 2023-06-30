package main

import (
	"github.com/OpenIMSDK/Open-IM-Server/internal/tools"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/cmd"
)

func main() {
	cronTaskCmd := cmd.NewCronTaskCmd()
	if err := cronTaskCmd.Exec(tools.StartCronTask); err != nil {
		panic(err.Error())
	}
}
