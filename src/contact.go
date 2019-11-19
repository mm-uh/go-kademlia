package kademlia

type KademliaContact struct {
	id   uint64
	ip   string
	port int
}

func (kc *KademliaContact) GetNodeId() uint64 {
	return kc.id
}

func (kc *KademliaContact) GetIP() string {
	return kc.ip
}

func (kc *KademliaContact) GetPort() int {
	return kc.port
}

func newKademliaContact(key uint64, ip string, port int) *KademliaContact {
	return &KademliaContact{
		id:   key,
		ip:   ip,
		port: port,
	}
}

func (kc *KademliaContact) PingRequestPingReply(key, sender_id, random_id string) {

}

func (kc *KademliaContact) FindNode(key, sender_id, looked_up_id, random_id string) {

}

func (kc *KademliaContact) FindValue(key, sender_id, looked_up_id, random_id string) {

}

type AddrPublishingNode struct {
	key  string
	port string
	ip   string
}

func (kc *KademliaContact) StoreMessage(key, sender_id string, addr AddrPublishingNode, random_id string) {

}

type ContactNode struct {
	id   string
	port string
	ip   string
}

func (kc *KademliaContact) NodeReply(key, sender_id, echoed_random_id string, addr []*ContactNode, random_id string) {

}

func (kc *KademliaContact) ValueReply(demultiplexer_key, sender_id, echoed_random_id, key string, addr []*ContactNode, random_id string) {

}
