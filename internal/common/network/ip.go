package network

import utils "github.com/OpenIMSDK/open_utils"

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
		return "0.0.0.0"
	} else {
		return configIP
	}
}
