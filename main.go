package main

import (
	kademlia "github.com/mm-uh/go-kademlia/src"
)

func main() {

	n := kademlia.NewLocalKademlia("localhost", 8080, 20, 3)
	n1 := kademlia.NewLocalKademlia("localhost", 8081, 20, 3)
	n1.JoinNetwork(n)
}
