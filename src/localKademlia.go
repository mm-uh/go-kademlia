package kademlia

type LocalKademlia struct {
	ft   *kademliaFingerTable
	ip   string
	port int
	id   KeyNode
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
}
