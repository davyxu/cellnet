package udppeer

import (
	"net"
)

type v4Address struct {
	IP   [4]byte
	Port int
}

type v6Address struct {
	IP   [16]byte
	Port int
}

type addressPair interface{}

func makeAddrKey(addr *net.UDPAddr) addressPair {

	switch len(addr.IP) {
	case net.IPv4len:
		var ret v4Address
		for i := 0; i < net.IPv4len; i++ {
			ret.IP[i] = addr.IP[i]
		}
		ret.Port = addr.Port

		return ret
	case net.IPv6len:
		var ret v6Address
		for i := 0; i < net.IPv6len; i++ {
			ret.IP[i] = addr.IP[i]
		}
		ret.Port = addr.Port

		return ret
	}

	return nil
}
