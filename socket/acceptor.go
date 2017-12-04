package socket

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/internal"
	"net"
	"sync"
)

// 接受器
type socketAcceptor struct {
	socketPeer
	internal.SessionManager

	// 保存侦听器
	l net.Listener

	// 侦听器的停止同步
	wg sync.WaitGroup
}

// 异步开始侦听
func (self *socketAcceptor) Start() cellnet.Peer {

	go self.listen(self.Address)

	return self
}

func (self *socketAcceptor) listen(address string) {

	// 侦听开始，添加1个任务
	self.wg.Add(1)

	// 在退出函数时，结束侦听任务
	defer self.wg.Done()

	var err error
	// 根据给定地址进行侦听
	self.l, err = net.Listen("tcp", address)

	// 如果侦听发生错误，打印错误并退出
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	log.Infof("#listen(%s) %s", self.Name(), self.Address)

	// 侦听循环
	for {

		// 新连接没有到来时，Accept是阻塞的
		conn, err := self.l.Accept()

		// 发生任何的侦听错误，打印错误并退出服务器
		if err != nil {
			break
		}

		go self.onNewSession(conn)
	}
}

func (self *socketAcceptor) onNewSession(conn net.Conn) {

	ses := newSession(conn, &self.socketPeer)

	ses.start()

}

// 停止侦听器
func (self *socketAcceptor) Stop() {
	self.l.Close()
	self.wg.Wait()
}

func init() {

	cellnet.RegisterPeerCreator("tcp.Acceptor", func(config cellnet.PeerConfig) cellnet.Peer {
		p := &socketAcceptor{
			SessionManager: internal.NewSessionManager(),
		}

		p.PeerConfig = config
		p.peerInterface = p

		return p
	})
}
