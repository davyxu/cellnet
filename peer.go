package cellnet

type Peer interface {
	Start(address string) Peer
	Stop()

	SetName(string)
	Name() string
}
