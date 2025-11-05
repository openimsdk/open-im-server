package main

import (
	"github.com/openimsdk/open-im-server/v3/pkg/common/cmd"
	"github.com/openimsdk/tools/system/program"
)

func main() {
	if err := cmd.NewMsgTransferCmd().Exec(); err != nil {
		program.ExitWithError(err)
	}
}
