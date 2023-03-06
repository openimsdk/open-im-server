package main

import (
	"OpenIM/internal/api"
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
		port, _ := cmd.Flags().GetInt("port")
		configFolderPath, _ := cmd.Flags().GetString("config_folder_path")
		fmt.Printf("Starting server on port %s with config file at %s\n", port, configFolderPath)
		if err := run(configFolderPath, port); err != nil {
			panic(err.Error())
		}
	},
}

func init() {
	startCmd.Flags().IntP("port", "p", 10002, "Port to listen on")
	startCmd.Flags().StringP("config_path", "c", "", "Path to config file folder")
}

func run(configFolderPath string, port int) error {
	if err := config.InitConfig(configFolderPath); err != nil {
		return err
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
	if err := startCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
