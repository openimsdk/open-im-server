package office

import (
	"github.com/OpenIMSDK/Open-IM-Server/internal/rpc/office"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/cmd"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
)

func main() {
	rpcCmd := cmd.NewRpcCmd("office")
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(); err != nil {
		panic(err.Error())
	}
	if err := rpcCmd.StartSvr(config.Config.RpcRegisterName.OpenImOfficeName, office.Start); err != nil {
		panic(err.Error())
	}
}
