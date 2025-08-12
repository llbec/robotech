package tcpclient

import (
	"net"
	"strconv"
)

func IsValidHost(s string) bool {
	host, portStr, err := net.SplitHostPort(s)
	if err != nil {
		return false
	}

	// 检查 IP 是否有效
	if ip := net.ParseIP(host); ip == nil {
		return false
	}

	// 检查端口是否是数字，且在合法范围
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		return false
	}

	return true
}
