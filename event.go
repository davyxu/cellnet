package cellnet

type SessionEvent interface {
	GetSession() CellID
}

type SessionPacket interface {
	SessionEvent

	GetPacket() *Packet
}
