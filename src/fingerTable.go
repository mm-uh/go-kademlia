package kademlia

import "fmt"

type kademliaFingerTable struct {
	kbuckets []KBucket
	id       KeyNode
}

func (ft *kademliaFingerTable) GetKBucket(index int) KBucket {
	fmt.Println(index)
	fmt.Println(len(ft.kbuckets))
	if index < 0 || index > len(ft.kbuckets) {
		return nil
	}

	fmt.Println(ft.kbuckets[index])
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
	fmt.Println(1)
	var index int = ft.id.Lenght() - 1
	fmt.Println(2)
	for ; index >= 0; index-- {
		if dist.IsActive(index) {
			break
		}
	}
	fmt.Println(3)
	kbucket := ft.GetKBucket(index)
	fmt.Println(4)
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
