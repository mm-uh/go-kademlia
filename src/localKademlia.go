package kademlia

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	serverRpc "github.com/mm-uh/rpc_udp/src/server"
	"github.com/sirupsen/logrus"

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
	time uint64
}

func NewLocalKademlia(ip string, port, k int, a int) *LocalKademlia {
	key := KeyNode(sha1.Sum([]byte(fmt.Sprintf("%s:%d", ip, port))))
	kBuckets := make([]KBucket, 0)
	for i := 0; i < 160; i++ {
		kBuckets = append(kBuckets, NewKademliaKBucket(k))
	}
	ft := &kademliaFingerTable{
		id:       &key,
		kbuckets: kBuckets,
	}
	return &LocalKademlia{
		ft:   ft,
		id:   &key,
		ip:   ip,
		port: port,
		k:    k,
		a:    a,
		time: 0,
	}
}

func (lk *LocalKademlia) GetContactInformation() *ContactInformation {
	return &ContactInformation{
		node: lk,
		time: lk.time,
	}
}
func (lk *LocalKademlia) JoinNetwork(node Kademlia) error {
	logrus.Info("Joining")
	lk.ft.Update(node)
	lk.nodeLookup(lk.GetNodeId())
	var index int = 0
	for ; index < lk.GetNodeId().Lenght(); index++ {
		kb, err := lk.ft.GetKBucket(index)
		if err != nil {
			return err
		}
		if len(kb.GetAllNodes()) != 0 {
			break
		}
	}

	index++
	for ; index < lk.GetNodeId().Lenght(); index++ {
		key := lk.ft.GetKeyFromKBucket(index)
		lk.nodeLookup(key)
	}
	logrus.Info("Joined")
	return nil
}

func (lk *LocalKademlia) Ping(ci *ContactInformation) bool {
	if ci != nil {
		lk.time = Max(lk.time, ci.time)
		lk.ft.Update(ci.node)
	}

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

func (lk *LocalKademlia) ClosestNodes(ci *ContactInformation, k int, id Key) ([]Kademlia, error) {
	if ci != nil {
		lk.time = Max(lk.time, ci.time)
		lk.ft.Update(ci.node)
	}
	return lk.ft.GetClosestNodes(k, id)
}

func (lk *LocalKademlia) Store(ci *ContactInformation, key Key, data interface{}) error {
	lk.time = Max(lk.time, ci.time)
	lk.ft.Update(ci.node)
	return lk.sm.Store(key, data)
}

func (lk *LocalKademlia) Get(ci *ContactInformation, id Key) (interface{}, error) {
	if ci != nil {
		lk.time = Max(lk.time, ci.time)
		lk.ft.Update(ci.node)
	}
	return lk.sm.Get(id)
}

func (lk *LocalKademlia) StoreOnNetwork(ci *ContactInformation, id Key, data interface{}) error {
	if ci != nil {
		lk.time = Max(lk.time, ci.time)
		lk.ft.Update(ci.node)
	}
	return nil
}

func (lk *LocalKademlia) GetFromNetwork(ci *ContactInformation, id Key) (interface{}, error) {
	if ci != nil {
		lk.time = Max(lk.time, ci.time)
		lk.ft.Update(ci.node)
	}
	return nil, nil
}

func (lk *LocalKademlia) GetInfo() string {
	ln := KBucketNodeInformation{
		Ip:   lk.GetIP(),
		Port: lk.GetPort(),
		Key:  fmt.Sprintf("%s", lk.GetNodeId()),
	}
	kbuckets := make(map[int][]KBucketNodeInformation, 0)
	for i := 0; i < lk.GetNodeId().Lenght(); i++ {
		kb, _ := lk.ft.GetKBucket(i)
		nodes := kb.GetAllNodes()
		if nodes != nil {
			for _, node := range nodes {
				n := KBucketNodeInformation{
					Ip:   node.GetIP(),
					Port: node.GetPort(),
					Key:  fmt.Sprintf("%s", node.GetNodeId()),
				}
				kbuckets[i] = append(kbuckets[i], n)
			}
		}
	}
	info := NodeInformation{
		KBucketNodeInformation: ln,
		Kbuckets:               kbuckets,
	}

	bytes, _ := json.Marshal(info)

	return string(bytes)
}

func (lk *LocalKademlia) nodeLookup(id Key) ([]Kademlia, error) {
	var round int = 1
	//Create structure to keep ordered nodes

	startNodes, err := lk.ft.GetClosestNodes(lk.a, id)
	if err != nil {
		return nil, err
	}
	if len(startNodes) == 0 {
		return nil, nil
	}
	//ToDo manage error
	dist, err := startNodes[0].GetNodeId().XOR(id)
	if err != nil {
		return nil, err
	}

	mu := sync.Mutex{}
	cond := sync.NewCond(&mu)

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
		go queryNode(node, id, round, lk.k, lk.GetContactInformation(), receivFromWorkers, cond)
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
				go queryNode(n.node, id, round, lk.k, lk.GetContactInformation(), receivFromWorkers, cond)
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
			cond.L.Lock()
			cond.Broadcast()
			cond.L.Unlock()
			return answ, nil
		}

	}

}

func (lk *LocalKademlia) RunServer(exited chan bool) {
	var h HandlerRPC
	h.kademlia = lk

	server := serverRpc.NewServer(h, lk.ip+":"+strconv.FormatInt(int64(lk.port), 10))

	// listen to incoming udp packets
	go server.ListenServer(exited)

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

func queryNode(node Kademlia, id Key, round, k int, ci *ContactInformation, send chan nodesPackage, finish *sync.Cond) {
	lookupFinished := make(chan bool)
	go func() {
		finish.L.Lock()
		finish.Wait()
		finish.L.Unlock()
		lookupFinished <- true
	}()
	nodes, _ := node.ClosestNodes(ci, k, id)
	channel := make(chan Kademlia)
	np := nodesPackage{
		round:       round,
		receivNodes: channel,
	}
	select {
	case send <- np:
		{
			for _, n := range nodes {
				channel <- n
			}
			close(channel)
		}
	case _ = <-lookupFinished:
		{
			return
		}
	}

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
