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

type Address struct {
	Scheme string
	Host   string
	Port   int
	Path   string
}

func (self *Address) String() string {
	if self.Scheme == "" {
		return fmt.Sprintf("%s:%d", self.Host, self.Port)
	}

	return fmt.Sprintf("%s://%s:%d%s", self.Scheme, self.Host, self.Port, self.Path)
}

func (self *Address) HostPort() string {

	return fmt.Sprintf("%s:%d", self.Host, self.Port)
}

// 在给定的端口范围内找到一个能用的端口 格式:
// scheme://host:minPort~maxPort/path
func DetectPort(addr string, fn func(*Address) (interface{}, error)) (interface{}, error) {

	var addrObj Address
	schemePos := strings.Index(addr, "://")

	// 移除scheme部分
	if schemePos != -1 {
		addrObj.Scheme = addr[:schemePos]
		addr = addr[schemePos+3:]
	}

	colonPos := strings.Index(addr, ":")

	if colonPos != -1 {
		addrObj.Host = addr[:colonPos]
	}

	addr = addr[colonPos+1:]

	rangePos := strings.Index(addr, "~")

	var minStr, maxStr string
	if rangePos != -1 {
		minStr = addr[:rangePos]

		slashPos := strings.Index(addr, "/")

		if slashPos != -1 {
			maxStr = addr[rangePos+1 : slashPos]
			addrObj.Path = addr[slashPos:]
		} else {
			maxStr = addr[rangePos:]
		}
	} else {
		slashPos := strings.Index(addr, "/")

		if slashPos != -1 {
			addrObj.Path = addr[slashPos:]
			minStr = addr[rangePos+1 : slashPos]
		} else {
			minStr = addr[rangePos+1:]
		}
	}

	// extract min port
	min, err := strconv.Atoi(minStr)
	if err != nil {
		return nil, ErrInvalidPortRange
	}

	var max int
	if maxStr != "" {
		// extract max port
		max, err = strconv.Atoi(maxStr)
		if err != nil {
			return nil, ErrInvalidPortRange
		}
	} else {
		max = min
	}

	for port := min; port <= max; port++ {

		addrObj.Port = port

		// 使用回调侦听
		ln, err := fn(&addrObj)
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

// 获取本地IP地址，有多重IP时，默认取第一个
func GetLocalIP() string {

	list, err := GetPrivateIPv4()
	if err != nil {
		return ""
	}

	if len(list) == 0 {
		return ""
	}

	return list[0].String()
}

// from consul

// GetPrivateIPv4 returns the list of private network IPv4 addresses on
// all active interfaces.
func GetPrivateIPv4() ([]*net.IPAddr, error) {
	addresses, err := activeInterfaceAddresses()
	if err != nil {
		return nil, fmt.Errorf("Failed to get interface addresses: %v", err)
	}

	var addrs []*net.IPAddr
	for _, rawAddr := range addresses {
		var ip net.IP
		switch addr := rawAddr.(type) {
		case *net.IPAddr:
			ip = addr.IP
		case *net.IPNet:
			ip = addr.IP
		default:
			continue
		}
		if ip.To4() == nil {
			continue
		}
		if !isPrivate(ip) {
			continue
		}
		addrs = append(addrs, &net.IPAddr{IP: ip})
	}
	return addrs, nil
}

// GetPublicIPv6 returns the list of all public IPv6 addresses
// on all active interfaces.
func GetPublicIPv6() ([]*net.IPAddr, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("Failed to get interface addresses: %v", err)
	}

	var addrs []*net.IPAddr
	for _, rawAddr := range addresses {
		var ip net.IP
		switch addr := rawAddr.(type) {
		case *net.IPAddr:
			ip = addr.IP
		case *net.IPNet:
			ip = addr.IP
		default:
			continue
		}
		if ip.To4() != nil {
			continue
		}
		if isPrivate(ip) {
			continue
		}
		addrs = append(addrs, &net.IPAddr{IP: ip})
	}
	return addrs, nil
}

// privateBlocks contains non-forwardable address blocks which are used
// for private networks. RFC 6890 provides an overview of special
// address blocks.
var privateBlocks = []*net.IPNet{
	parseCIDR("10.0.0.0/8"),     // RFC 1918 IPv4 private network address
	parseCIDR("100.64.0.0/10"),  // RFC 6598 IPv4 shared address space
	parseCIDR("127.0.0.0/8"),    // RFC 1122 IPv4 loopback address
	parseCIDR("169.254.0.0/16"), // RFC 3927 IPv4 link local address
	parseCIDR("172.16.0.0/12"),  // RFC 1918 IPv4 private network address
	parseCIDR("192.0.0.0/24"),   // RFC 6890 IPv4 IANA address
	parseCIDR("192.0.2.0/24"),   // RFC 5737 IPv4 documentation address
	parseCIDR("192.168.0.0/16"), // RFC 1918 IPv4 private network address
	parseCIDR("::1/128"),        // RFC 1884 IPv6 loopback address
	parseCIDR("fe80::/10"),      // RFC 4291 IPv6 link local addresses
	parseCIDR("fc00::/7"),       // RFC 4193 IPv6 unique local addresses
	parseCIDR("fec0::/10"),      // RFC 1884 IPv6 site-local addresses
	parseCIDR("2001:db8::/32"),  // RFC 3849 IPv6 documentation address
}

func parseCIDR(s string) *net.IPNet {
	_, block, err := net.ParseCIDR(s)
	if err != nil {
		panic(fmt.Sprintf("Bad CIDR %s: %s", s, err))
	}
	return block
}

func isPrivate(ip net.IP) bool {
	for _, priv := range privateBlocks {
		if priv.Contains(ip) {
			return true
		}
	}
	return false
}

// Returns addresses from interfaces that is up
func activeInterfaceAddresses() ([]net.Addr, error) {
	var upAddrs []net.Addr
	var loAddrs []net.Addr

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("Failed to get interfaces: %v", err)
	}

	for _, iface := range interfaces {
		// Require interface to be up
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addresses, err := iface.Addrs()
		if err != nil {
			return nil, fmt.Errorf("Failed to get interface addresses: %v", err)
		}

		if iface.Flags&net.FlagLoopback != 0 {
			loAddrs = append(loAddrs, addresses...)
			continue
		}

		upAddrs = append(upAddrs, addresses...)
	}

	if len(upAddrs) == 0 {
		return loAddrs, nil
	}

	return upAddrs, nil
}
