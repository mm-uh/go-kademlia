package kademlia

import (
	"encoding/json"
	"errors"
	"github.com/mm-uh/rpc_udp/src/util"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"time"
)

type Contacts []Kademlia

type KademliaContact struct {
	key  Key
	ip   string
	port int
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
	rpcBase := &util.RPCBase{
		MethodName: methodName,
	}
	response, err := kc.MakeRequest(rpcBase)
	if err != nil {
		return false
	}
	if response.Response != "Pong" {
		return false
	}
	return true
}

func (kc *KademliaContact) Store(key Key, i interface{}) error {
	methodName := "Store"
	rpcBase := &util.RPCBase{
		MethodName: methodName,
	}
	args := make([]interface{}, 0)
	args = append(args, key.GetString())
	args = append(args, i)
	rpcBase.Args = args
	response, err := kc.MakeRequest(rpcBase)
	if err != nil {
		return err
	}
	if response.Response != "true" {
		return errors.New("error storing value")
	}
	return nil
}

func (kc *KademliaContact) Get(key Key) (interface{}, error) {
	methodName := "Get"
	rpcBase := &util.RPCBase{
		MethodName: methodName,
	}
	args := make([]string, 0)
	args = append(args, key.GetString())
	rpcBase.Args = args

	response, err := kc.MakeRequest(rpcBase)
	if err != nil {
		return nil, err
	}
	var returnedKey Key
	err = returnedKey.GetFromString(response.Response)
	if err != nil {
		return nil, err
	}
	return returnedKey, nil
}

func (kc *KademliaContact) ClosestNodes(k int, key Key) ([]Kademlia, error) {
	methodName := "ClosesNodes"
	rpcBase := &util.RPCBase{
		MethodName: methodName,
	}
	args := make([]float64, 0)
	args = append(args, float64(k))
	rpcBase.Args = args

	response := Contacts{}
	nodeResponse, err := kc.MakeRequest(rpcBase)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(nodeResponse.Response), &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (kc *KademliaContact) MakeRequest(rpcBase *util.RPCBase) (*util.ResponseRPC, error) {

	service := kc.ip + ":" + strconv.Itoa(kc.port)

	RemoteAddr, err := net.ResolveUDPAddr("udp", service)

	conn, err := net.DialUDP("udp", nil, RemoteAddr)

	if err != nil {
		logrus.Warn(err)
		return nil, err
	}

	defer conn.Close()

	// write a message to server
	toSend, err := json.Marshal(rpcBase)
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

	timeout := time.Duration(1000)
	deadline := time.Now().Add(timeout)
	err = conn.SetReadDeadline(deadline)
	if err != nil {
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
