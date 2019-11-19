package kademlia

type KademliaNode struct {
	contact    *KademliaContact
	figerTable FingerTable
}
