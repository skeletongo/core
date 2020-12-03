package network

import (
	"github.com/skeletongo/core/module"
)

type IService interface {
	Start() error
	Update()
	Shutdown()
}

type Network struct {
	service map[int]IService
	close   bool
}

func (n *Network) Name() string {
	return "network"
}

func (n *Network) Init() {
	for i := 0; i < len(Config.Services); i++ {
		n.NewService(Config.Services[i])
	}
}

func (n *Network) Update() {
	for _, v := range n.service {
		v.Update()
	}
}

func (n *Network) Close() {
	if n.close {
		return
	}
	n.close = true

	if len(n.service) == 0 {
		module.Closed(n)
		return
	}

	for _, v := range n.service {
		go v.Shutdown()
	}
}

func (n *Network) ServiceClosed(cfg *SessionConfig) {
	delete(n.service, cfg.Id)
	if n.close && len(n.service) == 0 {
		module.Closed(n)
	}
}

func New() *Network {
	return &Network{
		service: make(map[int]IService),
	}
}

func (n *Network) NewService(sc *SessionConfig) IService {
	if n.close {
		return nil
	}

	var s IService
	if sc.IsClient {
		switch sc.Protocol {
		case "ws", "wss":

		case "udp":

		default:
			s = NewTCPClient(n, sc)
		}
	} else {
		switch sc.Protocol {
		case "ws", "wss":

		case "udp":

		default:
			s = NewTCPServer(n, sc)
		}
	}

	if s == nil {
		return nil
	}

	if err := s.Start(); err != nil {
		return nil
	}
	n.service[sc.Id] = s
	return s
}
