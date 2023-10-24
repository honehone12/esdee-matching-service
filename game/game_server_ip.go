package game

type GameServerIp interface {
	GetNextProcessIp() (string, string)
}

type DummyGameServerIp struct {
	address string
	port    string
}

func NewDummy(address string, port string) *DummyGameServerIp {
	return &DummyGameServerIp{
		address: address,
		port:    port,
	}
}

func (g *DummyGameServerIp) GetNextProcessIp() (string, string) {
	return g.address, g.port
}
