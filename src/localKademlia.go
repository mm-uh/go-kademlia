package kademlia

type LocalKademlia struct {
	ft   *kademliaFingerTable
	ip   string
	port int
	id   *KeyNode
	k    int
}

func NewLocalKademlia(ip string, port, k int) *LocalKademlia {
	id := KeyNodeFromSHA256()
	kBuckets := make([]KBucket, 0)
	for i := 0; i < 160; i++ {
		kBuckets = append(kBuckets, NewKademliaKBucket(k))
	}
	ft := &kademliaFingerTable{
		id:       *id,
		kbuckets: kBuckets,
	}
	return &LocalKademlia{
		ft:   ft,
		id:   id,
		ip:   ip,
		port: port,
		k:    k,
	}
}

func (lk *LocalKademlia) Ping() bool {
	return true
}

func (lk *LocalKademlia) GetIP() string {
	return lk.ip
}

func (lk *LocalKademlia) GetPort() int {
	return lk.port
}

func (lk *LocalKademlia) GetNodeID() Key {
	return lk.id
}

func (lk *LocalKademlia) ClosestNodes(k int, id Key) []Kademlia {
	return lk.ft.GetClosestNodes(k, id)
}

func (lk *LocalKademlia) Store(Key, data interface{}) error {
	return nil
}

func (lk *LocalKademlia) Get(id Key) (interface{}, error) {
	return nil, nil
}
