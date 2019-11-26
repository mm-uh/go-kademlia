package kademlia

import (
	"time"

	avl "github.com/mm-uh/go-avl/src"
)

type LocalKademlia struct {
	ft   *kademliaFingerTable
	ip   string
	port int
	id   *KeyNode
	k    int
	sm   StorageManager
	a    int
}

func NewLocalKademlia(ip string, port, k int, a int) *LocalKademlia {
	id := KeyNodeFromSHA256()
	kBuckets := make([]KBucket, 0)
	for i := 0; i < 160; i++ {
		kBuckets = append(kBuckets, NewKademliaKBucket(k))
	}
	ft := &kademliaFingerTable{
		id:       *id,
		kbuckets: kBuckets,
	}
	return &LocalKademlia{
		ft:   ft,
		id:   id,
		ip:   ip,
		port: port,
		k:    k,
		a:    a,
	}
}

func (lk *LocalKademlia) Ping() bool {
	return true
}

func (lk *LocalKademlia) GetIP() string {
	return lk.ip
}

func (lk *LocalKademlia) GetPort() int {
	return lk.port
}

func (lk *LocalKademlia) GetNodeID() Key {
	return lk.id
}

func (lk *LocalKademlia) ClosestNodes(k int, id Key) []Kademlia {
	return lk.ft.GetClosestNodes(k, id)
}

// func (lk *LocalKademlia) Store(Key, data interface{}) error {
// 	return nil
// }

// func (lk *LocalKademlia) Get(id Key) (interface{}, error) {
// 	return nil, nil
// }

func (lk *LocalKademlia) StoreOnNetwork(id Key, data interface{}) error {

	return nil
}

func (lk *LocalKademlia) GetFromNetwork(id Key) (interface{}, error) {
	return nil, nil
}

func (lk *LocalKademlia) nodeLookup(id Key) []Kademlia {
	var round int = 1
	//Create structure to keep ordered nodes

	startNodes := lk.ft.GetClosestNodes(lk.a, id)
	if len(startNodes) == 0 {
		return nil
	}
	//ToDo manage error
	dist, _ := startNodes[0].GetNodeId().XOR(id)
	Nodes, startNodes = avl.NewNode(dist, startNodes[0]), startNodes[1:]
	nextRoundMain := make(chan bool)
	nextRoundReceiver := make(chan bool)
	allNodesComplete := make(chan int)
	receivFromWorkers := make(chan nodesPackage)
	receivFromStorage := make(chan nodesPackage)
	go startRoundGuard(nextRoundMain, nextRoundReceiver, allNodesComplete)
	go replyReceiver(receivFromStorage, receivFromWorkers, nextRoundReceiver, allNodesComplete)
	for _, node := range startNodes {
		go queryNode(node, round, receivFromWorkers)
		//Add node to ordered structure
	}

	for {
		ended := <-nextRoundMain
		np := <-receivFromStorage
		nodes := np.receivNodes
		for node := range nodes {
			//add nodes to ordered structure
		}
		//get k first node from ordered structure
		//select a first not queried and repeat
		//if there is not queried nodes finish

	}

}

func startRoundGuard(nextRoundMain, nextRoundReceiver chan bool, allNodesComplete chan int) {
	var actRound int = 1
	for {
		timeout := make(chan bool)
		go func() { //timeout goroutine
			time.Sleep(1 * time.Second)
			timeout <- true
		}()

		for {
			select {
			case next := <-timeout:
				{

				}

			case round := <-allNodesComplete:
				{
					if round == actRound {
						break
					}
				}
			}
		}

		actRound++
		nextRoundMain <- true
		nextRoundReceiver <- true

	}

}

func queryNode(node Kademlia, round int, send chan nodesPackage) {

}

func replyReceiver(sendToMain, receivFromWorkers chan nodesPackage, nextRound chan bool, allNodesComplete chan int) {

}

type nodesPackage struct {
	round       int
	receivNodes chan Kademlia
}

type nodeLookup struct {
	queried bool
	node    Kademlia
}

type test struct{}

func (t test) Less(a test) (bool, error) {
	return true, nil
}
