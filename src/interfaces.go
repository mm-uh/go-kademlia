package kademlia

type Kademlia interface {
	Ping(*ContactInformation) bool
	Store(*ContactInformation, Key, interface{}) error
	Get(*ContactInformation, Key) (interface{}, error)
	StoreOnNetwork(*ContactInformation, Key, interface{}) error
	GetFromNetwork(*ContactInformation, Key) (interface{}, error)
	ClosestNodes(*ContactInformation, int, Key) ([]Kademlia, error)
	GetNodeId() Key
	GetIP() string
	GetPort() int
	JoinNetwork(Kademlia) error
}

type KBucket interface {
	Update(Kademlia)
	GetClosestNodes(int, Key) ([]Kademlia, error)
	GetAllNodes() []Kademlia
}

type FingerTable interface {
	GetClosestNodes(int, Key) ([]Kademlia, error)
	GetKBucket(int) (KBucket, error)
	Update(Kademlia) error
	GetKeyFromKBucket(k int) Key
}

type Key interface {
	XOR(other Key) (Key, error)
	IsActive(index int) bool
	Lenght() int
	Less(other interface{}) (bool, error)
	Equal(other interface{}) (bool, error)
	String() string
	GetFromString(string) error
}
