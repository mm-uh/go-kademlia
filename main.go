package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	kademlia "github.com/mm-uh/go-kademlia/src"
)

var Node *kademlia.LocalKademlia

func main() {
	ip := os.Args[1]
	portStr := os.Args[2]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("Invalid port")
	}

	gateway := len(os.Args) == 3

	ln := kademlia.NewLocalKademlia(ip, port, 20, 3)
	Node = ln
	exited := make(chan bool)
	ln.RunServer(exited)
	http.HandleFunc("/", EndpointHandler)
	go http.ListenAndServe(fmt.Sprintf(":%d", port+1000), nil)

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
	if s := <-exited; s {
		// Handle Error in method
		fmt.Println("We get an error listen server")
		return
	}
}

func EndpointHandler(w http.ResponseWriter, r *http.Request) {
	data := Node.GetInfo()
	fmt.Fprintf(w, "%s", data)
}
