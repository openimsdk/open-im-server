package utils

import (
	"errors"
	"fmt"
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
				fmt.Println(ipnet.IP.String())
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("no ip")
}
