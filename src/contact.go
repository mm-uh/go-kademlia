package kademlia

import (
	"encoding/json"
	"github.com/mm-uh/rpc_udp/src/util"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
)

type KademliaContact struct {
	key  Key
	ip   string
	port int
}

func getClient(addr string) Kademlia {
	return nil
}

func newKademliaContact(key Key, ip string, port int) *KademliaContact {
	return &KademliaContact{
		key:  key,
		ip:   ip,
		port: port,
	}
}

func (kc *KademliaContact) GetNodeId() Key {
	return kc.key
}

func (kc *KademliaContact) GetIP() string {
	return kc.ip
}

func (kc *KademliaContact) GetPort() int {
	return kc.port
}

func (kc *KademliaContact) Ping() bool {
	methodName := "Ping"
	rpcbase := &util.RPCBase{
		MethodName: methodName,
		Args:       make([]string, 0),
	}
	rpcbase.Args = append(rpcbase.Args, "Ping")
	response, err := kc.MakeRequest(rpcbase)
	if err != nil {
		return false
	}
	if response.Response != "Pong" {
		return false
	}
	return true
}

func (kc *KademliaContact) Store(key Key, i interface{}) error {
	// Do whatever todo for store some value
	return nil
}

func (kc *KademliaContact) Get(key Key) (interface{}, error) {
	return nil, nil
}

func (kc *KademliaContact) ClosesNodes(k int, key Key) []Kademlia {
	return nil
}

func (kc KademliaContact) MakeRequest(rpcbase *util.RPCBase) (*util.ResponseRPC, error) {

	service := kc.ip + ":" + strconv.Itoa(kc.port)

	RemoteAddr, err := net.ResolveUDPAddr("udp", service)

	conn, err := net.DialUDP("udp", nil, RemoteAddr)

	if err != nil {
		logrus.Warn(err)
		return nil, err
	}

	defer conn.Close()

	// write a message to server
	toSend, err := json.Marshal(rpcbase)
	if err != nil {
		logrus.Warn(err)
		return nil, err
	}

	message := []byte(string(toSend))

	_, err = conn.Write(message)

	if err != nil {
		logrus.Warn("Errorrr: " + err.Error())
		return nil, err
	}

	// receive message from server
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)

	var response util.ResponseRPC
	err = json.Unmarshal(buffer[:n], &response)
	if err != nil {
		logrus.Warn("Error Unmarshaling response")
		return nil, err
	}
	return &response, nil
}
