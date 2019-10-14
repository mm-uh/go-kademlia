package kademlia

type Kademlia interface{
	Ping() Contact
	Store(uint64, interface{}) error
	Get(uint64) (interface{}, error)
	ClosestNodes(int, uint64) []Contact
	GetContact() Contact
}


type KBucket interface{
	Update(Contact) error
	GetClosesNodes(int, uint64) []Contact
}

type Contact interface{
	GetNodeId() uint64
	GetIP() string
	GetPort() int
}

type FingerTable interface{
	GetClosestNodes(int, uint64) []Contact
	GetKBucket(int) KBucket
	GetKBucketOfContact(Contact) KBucket
}
