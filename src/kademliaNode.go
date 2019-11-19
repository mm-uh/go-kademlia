package kademlia

type KademliaNode struct {
	contact    *kademliaContact
	figerTable FingerTable
}
