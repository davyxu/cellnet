package nats

import (
	natstransport "github.com/davyxu/cellnet/transport/nats"
	"github.com/nats-io/nats.go"
	"time"
)

type MsgQueue struct {
	nc *nats.Conn
	// 重连间隔
	ReconnInterval time.Duration
	// Ping间隔
	PingInterval time.Duration

	OnSend func(msg interface{}) (payload []byte, err error)
	OnRecv func(payload []byte) (msg interface{}, err error)
}

func (self *MsgQueue) Publish(topic string, msg interface{}) error {

	payload, err := self.OnSend(msg)
	if err != nil {
		return err
	}

	return self.nc.Publish(topic, payload)
}

func (self *MsgQueue) Subscribe(topic string, callback func(msg interface{}, err error)) error {
	_, err := self.nc.Subscribe(topic, func(raw *nats.Msg) {
		msg, err := self.OnRecv(raw.Data)
		callback(msg, err)
	})

	return err
}

func (self *MsgQueue) Connect(addr string) error {
	nc, err := nats.Connect(addr,
		nats.ReconnectWait(self.ReconnInterval),
		nats.PingInterval(self.PingInterval))

	if err != nil {
		return err
	}

	self.nc = nc

	nc.Flush()

	return nil
}

func NewMsgQueue() *MsgQueue {
	self := &MsgQueue{
		OnRecv:         natstransport.RecvMessage,
		OnSend:         natstransport.SendMessage,
		ReconnInterval: time.Second * 5,
		PingInterval:   time.Second * 10,
	}

	return self
}
