package main

import (
	"fmt"
	kademlia "github.com/mm-uh/go-kademlia/src"
	"net"
	"strconv"
	"testing"
	"time"
)

func Test(t *testing.T) {
	ip := "localhost"
	var port int

	for i := 8080; i <= 10000; i++ {
		exist := availablePort(i)
		if exist {
			port = i
			break
		}
	}

	ln := kademlia.NewLocalKademlia(ip, port, 20, 3)
	Node = ln
	exited := make(chan bool)
	ln.RunServer(exited)

	time.Sleep(5*time.Second)

	for i := port + 1; i <= port+20; i++ {
		exist := availablePort(i)
		if exist {
			rn := kademlia.NewRemoteKademliaWithoutKey(ip, i)
			go ln.JoinNetwork(rn)
		}

	}

	if s := <-exited; s {
		// Handle Error in method
		fmt.Println("We get an error listen server")
		return
	}
}

func availablePort(i int) bool {
	ln, err := net.Listen("tcp", ":"+strconv.FormatInt(int64(i), 10))
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}
