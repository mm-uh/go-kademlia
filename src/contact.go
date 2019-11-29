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

type RemoteKademlia struct {
	key  Key
	ip   string
	port int
}

func NewKademliaContact(key Key, ip string, port int) *RemoteKademlia {
	return &RemoteKademlia{
		key:  key,
		ip:   ip,
		port: port,
	}
}

func (kc *RemoteKademlia) GetNodeId() Key {
	return kc.key
}

func (kc *RemoteKademlia) GetIP() string {
	return kc.ip
}

func (kc *RemoteKademlia) GetPort() int {
	return kc.port
}

func (kc *RemoteKademlia) Ping() bool {
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

func (kc *RemoteKademlia) StoreOnNetwork(key Key,i interface{}) error {
	methodName := "StoreOnNetwork"
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

func (kc *RemoteKademlia) GetFromNetwork(key Key) (interface{}, error) {
	methodName := "GetFromNetwork"
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

func (kc *RemoteKademlia) JoinNetwork(kademlia Kademlia) error {
	methodName := "JoinNetwork"
	rpcBase := &util.RPCBase{
		MethodName: methodName,
	}
	args := make([]interface{}, 0)
	args = append(args, kademlia.GetNodeId())
	args = append(args, kademlia.GetIP())
	args = append(args, kademlia.GetPort())
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

func (kc *RemoteKademlia) Store(key Key, i interface{}) error {
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

func (kc *RemoteKademlia) Get(key Key) (interface{}, error) {
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

//func (kc *RemoteKademlia) ClosestNodes(k int, key Key) ([]Kademlia, error) {
//	methodName := "ClosesNodes"
//	rpcBase := &util.RPCBase{
//		MethodName: methodName,
//	}
//	args := make([]float64, 0)
//	args = append(args, float64(k))
//	rpcBase.Args = args
//
//	response := Contacts{}
//	nodeResponse, err := kc.MakeRequest(rpcBase)
//	if err != nil {
//		return response, err
//	}
//	err = json.Unmarshal([]byte(nodeResponse.Response), &response)
//	if err != nil {
//		return response, err
//	}
//	return response, nil
//}

func (kc *RemoteKademlia) MakeRequest(rpcBase *util.RPCBase) (*util.ResponseRPC, error) {

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
		logrus.Warn("Error: " + err.Error())
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
