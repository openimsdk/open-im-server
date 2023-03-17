package network

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	utils "github.com/OpenIMSDK/open_utils"
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
