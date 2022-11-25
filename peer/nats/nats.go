package nats

import (
	"errors"
	natstransport "github.com/davyxu/cellnet/transport/nats"
	"github.com/nats-io/nats.go"
	"time"
)

var (
	ErrMsgqNotReady = errors.New("msgq not ready")
)

type MsgQueue struct {
	nc     *nats.Conn
	OnSend func(msg any) (payload []byte, err error)
	OnRecv func(payload []byte) (msg any, err error)
}

func (self *MsgQueue) Publish(topic string, msg any) error {

	payload, err := self.OnSend(msg)
	if err != nil {
		return err
	}

	if self.nc == nil {
		return ErrMsgqNotReady
	}

	return self.nc.Publish(topic, payload)
}

func (self *MsgQueue) Subscribe(topic string, callback func(msg any, err error) any) error {
	if self.nc == nil {
		return ErrMsgqNotReady
	}
	_, err := self.nc.Subscribe(topic, func(raw *nats.Msg) {
		msg, err := self.OnRecv(raw.Data)
		reply := callback(msg, err)
		if raw.Reply != "" && reply != nil {
			self.Publish(raw.Reply, reply)
		}
	})

	return err
}

func (self *MsgQueue) Request(topic string, msg any, timeout time.Duration) (resp any, retErr error) {

	payload, err := self.OnSend(msg)
	if err != nil {
		retErr = err
		return
	}

	if self.nc == nil {
		return nil, ErrMsgqNotReady
	}

	reply, err := self.nc.Request(topic, payload, timeout)
	if err != nil {
		retErr = err
		return
	}

	return self.OnRecv(reply.Data)
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
