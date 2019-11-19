package kademlia

type kademliaFingerTable struct {
	kbuckets []KBucket
}

func (ft *kademliaFingerTable) GetKBucket(index int) KBucket {
	if index < 0 || index > len(ft.kbuckets) {
		return nil
	}
	return ft.kbuckets[index]
}

func (ft *kademliaFingerTable) GetClosestNodes(int, uint64) []Contact {
	return nil
}
