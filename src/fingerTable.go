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

func (ft *kademliaFingerTable) GetClosestNodes(k int, key Key) []Kademlia {

	closestNodes := make([]Kademlia, 0)
	//ToDo Handle error
	dist, _ := ft.id.XOR(key)

	for i := dist.Lenght() - 1; i > 0; i-- {
		if dist.IsActive(i) {
			newClosestNodes := ft.kbuckets[i].GetClosestNodes(k-len(closestNodes), key)
			if newClosestNodes != nil {
				closestNodes = append(closestNodes, newClosestNodes...)
			}

		}
		if len(closestNodes) == k {
			return closestNodes
		}
	}

	for i := 0; i < dist.Lenght(); i++ {
		if !dist.IsActive(i) {
			newClosestNodes := ft.kbuckets[i].GetClosestNodes(k-len(closestNodes), key)
			if newClosestNodes != nil {
				closestNodes = append(closestNodes, newClosestNodes...)
			}
		}
		if len(closestNodes) == k {
			return closestNodes
		}
	}
	return closestNodes

}

func (ft *kademliaFingerTable) Update(node Kademlia) {

	dist, _ := ft.id.XOR(node.GetNodeId())
	var index int = ft.id.Lenght() - 1
	for ; index >= 0; index-- {
		if dist.IsActive(index) {
			break
		}
	}
	kbucket := ft.GetKBucket(index)
	kbucket.Update(node)
}

func (ft *kademliaFingerTable) GetKeyFromKBucket(k int) Key {
	myKey := ft.id
	artKey := KeyNode{}
	indexByte := k / 8
	indexBit := uint8(k % 8)
	for i, bt := range myKey {
		artKey[i] = bt
		if i == indexByte {
			artKey[i] = bt ^ 1<<(indexBit)
		}
	}

	return &artKey
}
