package network

import (
	utils "github.com/OpenIMSDK/open_utils"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
)

func GetRpcRegisterIP(configIP string) (string, error) {
	registerIP := configIP
	if registerIP == "" {
		ip, err := utils.GetLocalIP()
		if err != nil {
			return "", err
		}
		registerIP = ip
	}
	return registerIP, nil
}

func GetListenIP(configIP string) string {
	if configIP == "" {
		return constant.LocalHost
	} else {
		return configIP
	}
}
