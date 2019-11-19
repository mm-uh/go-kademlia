package kademlia

type Kademlia interface {
	Ping() Contact
	Store(Key, interface{}) error
	Get(Key) (interface{}, error)
	ClosestNodes(int, Key) []Contact
	GetContact() Contact
}

type KBucket interface {
	Update(Contact) error
	GetClosestNodes(int, Key) []Contact
	GetAllNodes() []Contact
}

type Contact interface {
	GetNodeId() uint64
	GetIP() string
	GetPort() int
	Ping() bool
}

type FingerTable interface {
	GetClosestNodes(int, Key) []Contact
	GetKBucket(int) KBucket
}

type Key interface {
	XOR(other Key) (Key, error)
	IsActive(index int) bool
	Lenght() int
}
