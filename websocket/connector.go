package websocket

import (
	"fmt"
	"strings"
	"time"

	"github.com/davyxu/cellnet"
	"golang.org/x/net/websocket"
)

type socketConnector struct {
	*peerBase
	*sessionMgr

	conn             *websocket.Conn
	autoReconnectSec int // 重连间隔时间, 0为不重连

	closeSignal chan bool

	working bool // 重入锁

	defaultSes cellnet.Session

	tryConnTimes int // 尝试连接次数
}

// 自动重连间隔=0不重连
func (self *socketConnector) SetAutoReconnectSec(sec int) {
	self.autoReconnectSec = sec
}

func (self *socketConnector) Start(address string) cellnet.Peer {
	args := strings.Split(address, ",")
	orign, protocol, url := "", "", ""
	if len(args) >= 1 {
		orign = args[0]
	}
	if len(args) >= 2 {
		url = args[1]
	}
	if len(args) >= 3 {
		protocol = args[2]
	}
	fmt.Println(args)
	if self.working {
		return self
	}

	go self.connect(url, protocol, orign)

	return self
}

const reportConnectFailedLimitTimes = 3

func (self *socketConnector) connect(address string, protocol string, orign string) {
	self.working = true

	for {

		self.tryConnTimes++

		// 开始连接
		cn, err := websocket.Dial(address, "", orign)

		// 连不上
		if err != nil {

			if self.tryConnTimes <= reportConnectFailedLimitTimes {
				log.Errorf("#connect failed(%s) %v", self.name, err.Error())
			}

			if self.tryConnTimes == reportConnectFailedLimitTimes {
				log.Errorf("(%s) continue reconnecting, but mute log", self.name)
			}

			// 没重连就退出
			if self.autoReconnectSec == 0 {
				break
			}

			// 有重连就等待
			time.Sleep(time.Duration(self.autoReconnectSec) * time.Second)

			// 继续连接
			continue
		}

		self.tryConnTimes = 0

		// 连上了, 记录连接
		self.conn = cn

		// 创建Session
		ses := newSession(NewPacketStream(cn), self, self)
		self.sessionMgr.Add(ses)
		self.defaultSes = ses

		log.Infof("#connected(%s) %s sid: %d", self.name, address, ses.id)

		// 内部断开回调
		ses.OnClose = func() {
			self.sessionMgr.Remove(ses)
			self.closeSignal <- true
		}

		// 抛出事件
		self.Post(self, NewSessionEvent(Event_SessionConnected, ses, nil))

		if <-self.closeSignal {

			self.conn = nil

			// 没重连就退出
			if self.autoReconnectSec == 0 {
				log.Errorln("退出")
				break
			}

			// 有重连就等待
			time.Sleep(time.Duration(self.autoReconnectSec) * time.Second)

			// 继续连接
			continue

		}

	}

	self.working = false
}

func (self *socketConnector) Stop() {

	if self.conn != nil {

		self.conn.Close()
	}

}

func (self *socketConnector) DefaultSession() cellnet.Session {
	return self.defaultSes
}

func NewConnector(evq cellnet.EventQueue) cellnet.Peer {
	self := &socketConnector{
		sessionMgr:  newSessionManager(),
		peerBase:    newPeerBase(evq),
		closeSignal: make(chan bool),
	}

	return self
}
