package utils

import (
	"Open_IM/pkg/common/config"
	"errors"
	"net"
)

var ServerIP = ""

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {

		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("no ip")
}

func GetRpcIP() (string, error) {
	registerIP := config.Config.RpcRegisterIP
	if registerIP == "" {
		ip, err := GetLocalIP()
		if err != nil {
			return "", err
		}
		registerIP = ip
	}
	return registerIP, nil
}
