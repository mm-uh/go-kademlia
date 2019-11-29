package kademlia

type Kademlia interface {
	Ping(*ContactInformation) bool
	Store(*ContactInformation, Key, interface{}) error
	Get(*ContactInformation, Key) (interface{}, error)
	StoreOnNetwork(*ContactInformation, Key, interface{}) error
	GetFromNetwork(*ContactInformation, Key) (interface{}, error)
	ClosestNodes(*ContactInformation, int, Key) []Kademlia
	GetNodeId() Key
	GetIP() string
	GetPort() int
	JoinNetwork(Kademlia)
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
	Equal(other interface{}) (bool, error)
}
