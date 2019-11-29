package kademlia

type Kademlia interface {
	Ping() bool
	Store(Key, interface{}) error
	Get(Key) (interface{}, error)
	StoreOnNetwork(Key, interface{}) error
	GetFromNetwork(Key) (interface{}, error)
	GetNodeId() Key
	GetIP() string
	GetPort() int
	JoinNetwork(Kademlia) error
}

type KBucket interface {
	Update(Kademlia)
	GetClosestNodes(int, Key) []Kademlia
	GetAllNodes() []Kademlia
}

type FingerTable interface {
	GetClosestNodes(int, Key) []Kademlia
	GetKBucket(int) KBucket
	Update(Kademlia)
	GetKeyFromKBucket(k int) Key
}

type Key interface {
	XOR(other Key) (Key, error)
	IsActive(index int) bool
	Lenght() int
	Less(other interface{}) (bool, error)
	GetString() string
	GetFromString(string) error
}
