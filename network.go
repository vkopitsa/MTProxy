package main

import (
	"sync"
)

type Network struct {
	servers          []string
	server_idle_cons map[int16][]*Server
	num_servers      int

	mux sync.Mutex
}

func NewNetwork(servers []string) *Network {
	n := &Network{
		servers:     servers,
		num_servers: len(servers),
	}
	n.InitServers()
	return n
}

func (n *Network) InitServers() {
	n.server_idle_cons = make(map[int16][]*Server)

	for i, s := range n.servers {
		for ii := 0; ii < 4; ii++ {
			idDc := int16(i)

			server := NewServer(s, idDc)

			err := server.Run()
			if err != nil {
				continue
			}

			if _, ok := n.server_idle_cons[idDc]; !ok {
				n.server_idle_cons[idDc] = make([]*Server, 0, 4)
			}
			n.server_idle_cons[idDc] = append(n.server_idle_cons[idDc], server)
		}
	}
}

func (n *Network) GetServer(idDc int16) *Server {
	n.mux.Lock()
	s := n.server_idle_cons[idDc][0]
	n.server_idle_cons[idDc] = n.server_idle_cons[idDc][1:]
	n.checkServers(idDc)
	n.mux.Unlock()

	return s
}

func (n *Network) checkServers(idDc int16) {
	ssfdf := 3 - len(n.server_idle_cons[idDc])
	if ssfdf <= 0 {
		return
	}
	for i := 0; i < ssfdf; i++ {
		s := n.servers[idDc]
		server := NewServer(s, idDc)

		n.server_idle_cons[idDc] = append(n.server_idle_cons[idDc], server)

		server.Run()
	}
}
