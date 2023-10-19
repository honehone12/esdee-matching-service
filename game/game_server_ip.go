package game

import "fmt"

type GameServerIp interface {
	GetNextProcessIp() string
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

func (g *DummyGameServerIp) GetNextProcessIp() string {
	return fmt.Sprintf("%s:%s", g.address, g.port)
}
