package kademlia

type kademliaContact struct{
	id	uint64
	ip	string
	port	int
}

func (kc *kademliaContact) GetNodeId() uint64{
	return kc.id
}

func (kc *kademliaContact) GetIP() string{
	return kc.ip
}

func (kc *kademliaContact) GetPort() int{
	return kc.port
}

func newKademliaContact(key uint64, ip string, port int) *kademliaContact{
	return &kademliaContact{
		id: key,
		ip: ip,
		port: port,
	}
}
