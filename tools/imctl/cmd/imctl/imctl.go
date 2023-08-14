// iamctl is the command line tool for iam platform.
package main

import (
	"os"

	"github.com/OpenIMSDK/Open-IM-Server/tools/imctl/internal/imctl/cmd"
)

func main() {
	command := cmd.NewDefaultIMCtlCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
