package kademlia

type KademliaContact struct {
	id   uint64
	ip   string
	port int
}

func newKademliaContact(key uint64, ip string, port int) *KademliaContact {
	return &KademliaContact{
		id:   key,
		ip:   ip,
		port: port,
	}
}

func (kc *KademliaContact) GetNodeId() uint64 {
	return kc.id
}

func (kc *KademliaContact) GetIP() string {
	return kc.ip
}

func (kc *KademliaContact) GetPort() int {
	return kc.port
}

func (kc *KademliaContact) Ping() int {
	return kc.port
}

func (kc *KademliaContact) Store(id uint64, i interface{}) error {
	return nil
}

func (kc *KademliaContact) Get(id uint64) (interface{}, error) {
	return nil, nil
}

func (kc *KademliaContact) ClosesNodes(k int, id uint64) []Contact {
	return nil
}

func (kc *KademliaContact) GetContact() Contact {
	return nil
}
