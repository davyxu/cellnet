package util

import (
	"errors"
	"fmt"
	"github.com/davyxu/cellnet"
	"net"
	"strconv"
	"strings"
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

var (
	ErrInvalidPortRange = errors.New("invalid port range")
)

// 在给定的端口范围内找到一个能用的端口 格式: localhost:5000~6000
func DetectPort(addr string, fn func(string) (net.Listener, error)) (net.Listener, error) {
	// host:port 或 host:min~max
	parts := strings.Split(addr, ":")

	// host:port格式
	if len(parts) < 2 {
		return fn(addr)
	}

	// 间隔分割
	ports := strings.Split(parts[len(parts)-1], "~")

	// 单独的端口
	if len(ports) < 2 {
		return fn(addr)
	}

	// extract min port
	min, err := strconv.Atoi(ports[0])
	if err != nil {
		return nil, ErrInvalidPortRange
	}

	// extract max port
	max, err := strconv.Atoi(ports[1])
	if err != nil {
		return nil, ErrInvalidPortRange
	}

	host := parts[0]

	for port := min; port <= max; port++ {

		// 使用回调侦听
		ln, err := fn(fmt.Sprintf("%s:%d", host, port))
		if err == nil {
			return ln, nil
		}

		// hit max port
		if port == max {
			return nil, err
		}
	}

	return nil, fmt.Errorf("unable to bind to %s", addr)
}
