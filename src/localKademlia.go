package kademlia

import (
	"crypto/sha1"
	"fmt"
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
	key := KeyNode(sha1.Sum([]byte(fmt.Sprintf("%s:%d", ip, port))))
	kBuckets := make([]KBucket, 0)
	for i := 0; i < 160; i++ {
		kBuckets = append(kBuckets, NewKademliaKBucket(k))
	}
	ft := &kademliaFingerTable{
		id:       key,
		kbuckets: kBuckets,
	}
	return &LocalKademlia{
		ft:   ft,
		id:   &key,
		ip:   ip,
		port: port,
		k:    k,
		a:    a,
	}
}

func (lk *LocalKademlia) JoinNetwork(node Kademlia) {
	lk.ft.Update(node)
	lk.nodeLookup(lk.GetNodeId())
	var index int = 0
	for ; index < lk.GetNodeId().Lenght(); index++ {
		if len(lk.ft.GetKBucket(index).GetAllNodes()) != 0 {
			break
		}
	}

	index++
	for ; index < lk.GetNodeId().Lenght(); index++ {
		key := lk.ft.GetKeyFromKBucket(index)
		lk.nodeLookup(key)
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

func (lk *LocalKademlia) GetNodeId() Key {
	return lk.id
}

func (lk *LocalKademlia) ClosestNodes(k int, id Key) []Kademlia {
	return lk.ft.GetClosestNodes(k, id)
}

func (lk *LocalKademlia) Store(key Key, data interface{}) error {
	return lk.sm.Store(key, data)
}

func (lk *LocalKademlia) Get(id Key) (interface{}, error) {
	return lk.sm.Get(id)
}

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
	Nodes, startNodes := avl.NewNode(dist, newNodeLookup(startNodes[0])), startNodes[1:]
	node, _ := Nodes.Value.(nodeLookup)
	node.queried = true
	nextRoundMain := make(chan bool)
	nextRoundReceiver := make(chan bool)
	allNodesComplete := make(chan int)
	receivFromWorkers := make(chan nodesPackage)
	receivFromStorage := make(chan nodesPackage)
	endStorage := make(chan bool)
	endGuard := make(chan bool)
	go startRoundGuard(nextRoundMain, nextRoundReceiver, endGuard, allNodesComplete)
	go replyReceiver(receivFromStorage, receivFromWorkers, nextRoundReceiver, endStorage, allNodesComplete, lk.a)
	for _, node := range startNodes {
		go queryNode(node, id, round, lk.k, receivFromWorkers)
		//Add node to ordered structure
		dist, _ := node.GetNodeId().XOR(id)
		nl := newNodeLookup(node)
		nl.queried = true
		newNode := avl.NewNode(dist, nl)
		Nodes = avl.Insert(Nodes, newNode)
	}

	for {
		_ = <-nextRoundMain
		round++
		np := <-receivFromStorage
		nodes := np.receivNodes
		for node := range nodes {
			//add nodes to ordered structure
			dist, _ := node.GetNodeId().XOR(id)
			newNode := avl.NewNode(dist, newNodeLookup(node))
			Nodes = avl.Insert(Nodes, newNode)
		}
		mins := Nodes.GetKMins(lk.k)
		asked := 0
		for _, node := range mins {
			n, ok := node.Value.(*nodeLookup)
			if !ok {
				panic("Incorrect type")
			}
			if !n.queried {
				n.queried = true
				go queryNode(n.node, id, round, lk.k, receivFromWorkers)
				asked++
			}
			if asked == lk.a {
				break
			}
		}
		if asked == 0 {
			answ := make([]Kademlia, 0)
			for _, node := range Nodes.GetKMins(lk.k) {
				n, ok := node.Value.(*nodeLookup)
				if !ok {
					panic("Incorrect type")
				}
				answ = append(answ, n.node)
			}
			endGuard <- true
			endStorage <- true
			return answ
		}

	}

}

func startRoundGuard(nextRoundMain, nextRoundReceiver, lookupEnd chan bool, allNodesComplete chan int) {
	var actRound int = 1
	for {
		timeout := make(chan bool)
		go func() { //timeout goroutine
			time.Sleep(1 * time.Second)
			timeout <- true
		}()

		select {
		case _ = <-timeout:
			{
			}
		case round := <-allNodesComplete:
			{
				if round == actRound {
					break
				}
			}
		case _ = <-lookupEnd:
			{
				return
			}
		}

		actRound++
		nextRoundMain <- true
		nextRoundReceiver <- true

	}

}

func queryNode(node Kademlia, id Key, round, k int, send chan nodesPackage) {
	nodes := node.ClosestNodes(k, id)
	channel := make(chan Kademlia)
	np := nodesPackage{
		round:       round,
		receivNodes: channel,
	}
	send <- np
	for _, n := range nodes {
		channel <- n
	}
	close(channel)
}

func replyReceiver(sendToMain, receivFromWorkers chan nodesPackage, nextRound, lookupEnd chan bool, allNodesComplete chan int, a int) {
	var actRound int = 1
	var finished int = 0
	var nodesForSend []Kademlia
	for {

		select {
		case np := <-receivFromWorkers:
			{
				channel := np.receivNodes
				for node := range channel {
					nodesForSend = append(nodesForSend, node)
				}
				if np.round == actRound {
					finished++
				}
				if finished == a {
					allNodesComplete <- actRound
				}
			}

		case _ = <-nextRound:
			{
				actRound++
				channel := make(chan Kademlia)
				np := nodesPackage{
					round:       actRound,
					receivNodes: channel,
				}
				sendToMain <- np
				for _, node := range nodesForSend {
					channel <- node
				}
				close(channel)
				nodesForSend = make([]Kademlia, 0)
			}

		case _ = <-lookupEnd:
			{
				return
			}
		}

	}
}

type nodesPackage struct {
	round       int
	receivNodes chan Kademlia
}

type nodeLookup struct {
	queried bool
	node    Kademlia
}

func newNodeLookup(node Kademlia) *nodeLookup {
	return &nodeLookup{
		queried: false,
		node:    node,
	}
}
