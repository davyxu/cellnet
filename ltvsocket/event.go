package ltvsocket

import (
	"github.com/davyxu/cellnet"
)

type EventClose struct {
	error
}

type EventConnectError struct {
	error
}

type EventListenError struct {
	error
}

type EventNewSession interface {
	Stream() cellnet.IPacketStream
}

type EventConnected struct {
	stream cellnet.IPacketStream
}

func (self EventConnected) Stream() cellnet.IPacketStream {
	return self.stream
}

type EventAccepted struct {
	stream cellnet.IPacketStream
}

func (self EventAccepted) Stream() cellnet.IPacketStream {
	return self.stream
}
