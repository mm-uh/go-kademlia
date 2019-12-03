package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
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

	ln := kademlia.NewLocalKademlia(ip, port, 5, 3)
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

	fmt.Println("Get or Save??")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		switch scanner.Text() {
		case "save":
			{
				fmt.Println("Insert Key")
				scanner.Scan()
				keyStr := scanner.Text()
				key := kademlia.KeyNode{}
				hash := sha1.Sum([]byte(keyStr))
				key.GetFromString(hex.EncodeToString(hash[:]))
				fmt.Println("Insert Value")
				scanner.Scan()
				val := scanner.Text()
				//_, _ = ln.GetAndLock(ln.GetContactInformation(), &key)
				err := ln.GetLock(ln.GetContactInformation(), &key)
				if err != nil {
					fmt.Println("COULD NOT UPDATE")
					continue
				}
				ln.StoreOnNetwork(ln.GetContactInformation(), &key, val)
				ln.LeaveLock(ln.GetContactInformation(), &key)
				//ln.StoreAndUnlock(ln.GetContactInformation(), &key, val)
			}

		case "get":
			{
				fmt.Println("Insert Key")
				scanner.Scan()
				keyStr := scanner.Text()
				key := kademlia.KeyNode{}
				hash := sha1.Sum([]byte(keyStr))
				key.GetFromString(hex.EncodeToString(hash[:]))
				//val, err := ln.GetAndLock(ln.GetContactInformation(), &key)
				//ln.StoreAndUnlock(ln.GetContactInformation(), &key, val)
				err := ln.GetLock(ln.GetContactInformation(), &key)
				if err != nil {
					fmt.Println("COULD NOT GET VALUE")
					continue
				}
				val, err := ln.GetFromNetwork(ln.GetContactInformation(), &key)
				ln.LeaveLock(ln.GetContactInformation(), &key)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println(val)
				}

			}
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
