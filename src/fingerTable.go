package kademlia

type kademliaFingerTable struct {
	kbuckets []KBucket
	id       uint64
}

func (ft *kademliaFingerTable) GetKBucket(index int) KBucket {
	if index < 0 || index > len(ft.kbuckets) {
		return nil
	}
	return ft.kbuckets[index]
}

func (ft *kademliaFingerTable) GetClosestNodes(k int, key uint64) []Contact {
	remain := k
	dist := XOR(key, ft.id)
	for remain > 0 {

	}
	return nil
}
