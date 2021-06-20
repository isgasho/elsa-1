package utils

import (
	"net"
)

const (
	LocalIp   = "127.0.0.1"
	LocalHost = "localhost"
)

// get local ip address
func GetLocalIp() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addresses {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return "127.0.0.1"
}
