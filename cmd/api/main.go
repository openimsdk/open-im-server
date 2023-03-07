package main

import (
	"OpenIM/internal/api"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/log"
	"fmt"
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
		configFolderPath, _ := cmd.Flags().GetString(constant.FlagConf)
		if err := run(configFolderPath, port); err != nil {
			panic(err.Error())
		}
	},
}

func init() {
	startCmd.Flags().IntP(constant.FlagPort, "p", 0, "Port to listen on")
	startCmd.Flags().StringP(constant.FlagConf, "c", "", "Path to config file folder")
}

func run(configFolderPath string, port int) error {
	if err := config.InitConfig(configFolderPath); err != nil {
		return err
	}
	if port == 0 {
		port = config.Config.Api.GinPort[0]
	}
	log.NewPrivateLog(constant.LogFileName)
	router := api.NewGinRouter()
	address := constant.LocalHost + ":" + strconv.Itoa(port)
	if config.Config.Api.ListenIP != "" {
		address = config.Config.Api.ListenIP + ":" + strconv.Itoa(port)
	}
	fmt.Println("start api server, address: ", address, ", OpenIM version: ", constant.CurrentVersion)
	err := router.Run(address)
	if err != nil {
		log.Error("", "api run failed ", address, err.Error())
		return err
	}
	return nil
}

func main() {
	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(startCmd)
	if err := startCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
