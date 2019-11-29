package main

import (
	"fmt"

	kademlia "github.com/mm-uh/go-kademlia/src"
)

func main() {

	n := kademlia.NewLocalKademlia("localhost", 8080, 20, 3)
	n1 := kademlia.NewLocalKademlia("localhost", 8081, 20, 3)
	n2 := kademlia.NewLocalKademlia("localhost", 8082, 20, 3)
	fmt.Println(n.GetNodeId())
	fmt.Println(n1.GetNodeId())
	dist, _ := n.GetNodeId().XOR(n1.GetNodeId())
	var index int = n.GetNodeId().Lenght() - 1
	for {
		if index < 0 {
			break
		}
		if dist.IsActive(index) {
			break
		}
		index--
	}
	fmt.Println(dist)
	fmt.Println(index)

	n1.JoinNetwork(n)
	n2.JoinNetwork(n)
	test := n2.ClosestNodes(20, n.GetNodeId())
	fmt.Println(len(test))
}
