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
	dist := ft.id.XOR(key)
	actIndex := dist.Lenght() - 1
	for len(closestNodes) < k {

		for !dist.IsActive(actIndex) && actIndex > 0 {
			actIndex--
		}
		closestNodes = append(closestNodes, ft.kbuckets[actIndex].GetClosestNodes(k-len(closestNodes), key)...)

	}

	return nil
}
