package main

import (
	"OpenIM/internal/api"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/mw"
	"fmt"
	"github.com/OpenIMSDK/openKeeper"
	"github.com/spf13/cobra"
	"os"
	"strconv"

	"OpenIM/pkg/common/constant"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt(constant.FlagPort)
		if err := run(port); err != nil {
			panic(err.Error())
		}
	},
}

func main() {
	startCmd.Flags().IntP(constant.FlagPort, "p", 0, "Port to listen on")
	startCmd.Flags().StringP(constant.FlagConf, "c", "", "Path to config file folder")
	rootCmd := cmd.NewRootCmd()
	rootCmd.Command.AddCommand(startCmd)
	if err := startCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(port int) error {
	if port == 0 {
		port = config.Config.Api.GinPort[0]
	}
	zk, err := openKeeper.NewClient(nil, "", 10, "", "")
	if err != nil {
		return err
	}
	log.NewPrivateLog(constant.LogFileName)
	zk.AddOption(mw.GrpcClient())
	router := api.NewGinRouter(zk)
	address := constant.LocalHost + ":" + strconv.Itoa(port)
	if config.Config.Api.ListenIP != "" {
		address = config.Config.Api.ListenIP + ":" + strconv.Itoa(port)
	}
	fmt.Println("start api server, address: ", address, ", OpenIM version: ", constant.CurrentVersion)
	err = router.Run(address)
	if err != nil {
		log.Error("", "api run failed ", address, err.Error())
		return err
	}
	return nil
}
