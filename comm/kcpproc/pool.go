package kcpproc

import (
	"sync"
	"time"
)

const (

	// maximum packet size
	mtuLimit = 1500
)

var (
	// global packet buffer
	// shared among sending/receiving/FEC
	xmitBuf sync.Pool
)

func init() {
	xmitBuf.New = func() interface{} {
		return make([]byte, mtuLimit)
	}
}

func currentMs() uint32 {
	return uint32(time.Now().UnixNano() / int64(time.Millisecond))
}
