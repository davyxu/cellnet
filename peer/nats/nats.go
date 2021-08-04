package nats

import (
	natstransport "github.com/davyxu/cellnet/transport/nats"
	"github.com/nats-io/nats.go"
)

type MsgQueue struct {
	nc     *nats.Conn
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

func (self *MsgQueue) Connect(addr string, options ...nats.Option) error {
	nc, err := nats.Connect(addr, options...)

	if err != nil {
		return err
	}

	self.nc = nc

	nc.Flush()

	return nil
}

func NewMsgQueue() *MsgQueue {
	self := &MsgQueue{
		OnRecv: natstransport.RecvMessage,
		OnSend: natstransport.SendMessage,
	}

	return self
}
