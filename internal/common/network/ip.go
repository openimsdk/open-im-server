package network

import utils "github.com/OpenIMSDK/open_utils"

func GetRpcIP(configIP string) (string, error) {
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
