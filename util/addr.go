package util

import (
	"fmt"
	"net"
	"strconv"
)

func SpliteAddress(addr string) (host string, port int, err error) {

	var portStr string

	host, portStr, err = net.SplitHostPort(addr)

	if err != nil {
		return "", 0, err
	}

	port, err = strconv.Atoi(portStr)

	if err != nil {
		return "", 0, err
	}

	return
}

func JoinAddress(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}
