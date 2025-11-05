package main

import (
	"github.com/openimsdk/open-im-server/v3/pkg/common/cmd"
	"github.com/openimsdk/tools/system/program"
)

func main() {
	if err := cmd.NewMsgRpcCmd().Exec(); err != nil {
		program.ExitWithError(err)
	}
}
