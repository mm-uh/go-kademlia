package kademlia

type kademliaFingerTable struct {
	kbuckets []KBucket
	id       KeyNode
}

func (ft *kademliaFingerTable) GetKBucket(index int) KBucket {
	if index < 0 || index > len(ft.kbuckets) {
		return nil
	}
	return ft.kbuckets[index]
}

func (ft *kademliaFingerTable) GetClosestNodes(k int, key *KeyNode) []Contact {

	closestNodes := make([]Contact, 0)
	//ToDo Handle error
	dist, _ := ft.id.XOR(key)

	for i := dist.Lenght() - 1; i > 0; i-- {
		if dist.IsActive(i) {
			closestNodes = append(closestNodes, ft.kbuckets[i].GetClosestNodes(k-len(closestNodes), key)...)
		}
		if len(closestNodes) == k {
			return closestNodes
		}
	}

	for i := 0; i < dist.Lenght(); i++ {
		if !dist.IsActive(i) {
			closestNodes = append(closestNodes, ft.kbuckets[i].GetClosestNodes(k-len(closestNodes), key)...)
		}
		if len(closestNodes) == k {
			return closestNodes
		}
	}

	return closestNodes

}
