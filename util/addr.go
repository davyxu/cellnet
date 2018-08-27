package util

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"net"
	"strconv"
)

// 将地址拆分为ip和端口
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

// 将ip和端口合并为地址
func JoinAddress(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// 获取session远程的地址
func GetRemoteAddrss(ses cellnet.Session) (string, bool) {
	if c, ok := ses.Raw().(net.Conn); ok {
		return c.RemoteAddr().String(), true
	}

	return "", false
}
