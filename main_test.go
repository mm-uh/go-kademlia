package main

import (
	"fmt"
	kademlia "github.com/mm-uh/go-kademlia/src"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"testing"
	"time"
)

func Test(t *testing.T) {
	logrus.SetLevel(logrus.InfoLevel)
	ip := "localhost"
	var port int
	timer := make(chan bool)
	for i := 8080; i <= 10000; i++ {
		exist := availablePort(i)
		if exist {
			port = i
			break
		}
	}
	go Wait(timer, 30)
	ln := kademlia.NewLocalKademlia(ip, port, 20, 3)
	Node = ln
	exited := make(chan bool)
	ln.RunServer(exited)

	time.Sleep(5 * time.Second)

	for i := port + 1; i <= port+20; i++ {
		exist := availablePort(i)
		if exist {
			rn := kademlia.NewRemoteKademliaWithoutKey(ip, i)
			go ln.JoinNetwork(rn)
		}

	}
	select {
	case <-timer:
		fmt.Println("Finish")
		return
	case <-exited:
		fmt.Println("We get an error listen server")
		return

	}
}

func Wait(timer chan bool, i int) {
	time.Sleep((time.Duration(i)) * time.Second)
	timer <- true
}

func availablePort(i int) bool {
	ln, err := net.Listen("tcp", ":"+strconv.FormatInt(int64(i), 10))
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}
