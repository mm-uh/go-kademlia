package main

import (
	"os"
	"strconv"

	kademlia "github.com/mm-uh/go-kademlia/src"
)

func main() {
	ip := os.Args[1]
	portStr := os.Args[2]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("Invalid port")
	}

	gateway := len(os.Args) == 3

	ln := kademlia.NewLocalKademlia(ip, port, 20, 3)
	ln.RunServer()

	if !gateway {
		ipForJoin := os.Args[3]
		portForJoinStr := os.Args[4]
		portForJoin, err := strconv.Atoi(portForJoinStr)
		if err != nil {
			panic("Invalid port for join")
		}
		rn := kademlia.NewRemoteKademliaWithoutKey(ipForJoin, portForJoin)
		err = ln.JoinNetwork(rn)
		if err != nil {
			panic("Can't Join")
		}
	}
}
