package utils

import (
	"Open_IM/src/common/config"
	"net"
)

var ServerIP = ""

func init() {
	//fixme In the configuration file, ip takes precedence, if not, get the valid network card ip of the machine
	if config.Config.ServerIP != "" {
		ServerIP = config.Config.ServerIP
		return
	}
	//fixme Get the ip of the local network card
	netInterfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(netInterfaces); i++ {
		//Exclude useless network cards by judging the net.flag Up flag
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			address, _ := netInterfaces[i].Addrs()
			for _, addr := range address {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						ServerIP = ipNet.IP.String()
						return
					}
				}
			}
		}
	}
}
