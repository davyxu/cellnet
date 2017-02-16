package socket

import "github.com/davyxu/cellnet"

// Peer默认的编码
var PeerDefaultCodec string = "pb"

// 系统事件默认都是pb
var sysEventCodec = cellnet.FetchCodec("pb")
