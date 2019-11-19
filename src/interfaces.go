package kademlia

type Kademlia interface {
	Ping() Contact
	Store(uint64, interface{}) error
	Get(uint64) (interface{}, error)
	ClosestNodes(int, uint64) []Contact
	GetContact() Contact
}

type KBucket interface {
	Update(Contact) error
	GetClosestNodes(int, uint64) []Contact
	GetAllNodes() []Contact
}

type Contact interface {
	GetNodeId() uint64
	GetIP() string
	GetPort() int
	Ping() bool
}

type FingerTable interface {
	GetClosestNodes(int, uint64) []Contact
	GetKBucket(int) KBucket
}
