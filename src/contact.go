package kademlia

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mm-uh/rpc_udp/src/util"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"strings"

	"time"
)

type Contacts []Kademlia

type RemoteKademlia struct {
	key  Key
	ip   string
	port int
}

func NewRemoteKademlia(key Key, ip string, port int) *RemoteKademlia {
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

func (kc *RemoteKademlia) Ping(info *ContactInformation) bool {
	methodName := "Ping"
	rpcBase := &util.RPCBase{
		MethodName: methodName,
	}
	args := make([]string, 0)
	args = append(args, contactInfoToString(info))
	rpcBase.Args = args

	response, err := kc.MakeRequest(rpcBase)
	if err != nil {
		return false
	}
	if response.Response != "Pong" {
		return false
	}
	return true
}

func (kc *RemoteKademlia) StoreOnNetwork(info *ContactInformation, key Key,i interface{}) error {
	methodName := "StoreOnNetwork"
	rpcBase := &util.RPCBase{
		MethodName: methodName,
	}
	args := make([]interface{}, 0)
	args = append(args, contactInfoToString(info))
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

func (kc *RemoteKademlia) ClosestNodes(cInfo *ContactInformation, k int,key Key) ([]Kademlia, error) {
	methodName := "ClosestNodes"
	rpcBase := &util.RPCBase{
		MethodName: methodName,
	}
	args := make([]interface{}, 0)
	args = append(args, contactInfoToString(cInfo))
	args = append(args, strconv.FormatInt(int64(k), 10))
	args = append(args, key.GetString())
	rpcBase.Args = args

	nodeResponse, err := kc.MakeRequest(rpcBase)
	if err != nil {
		return nil, err
	}
	
	returnedNodes, err := getKademliaNodes(nodeResponse.Response)
	if err != nil {
		return nil, err
	}
	return returnedNodes, nil
}

func (kc *RemoteKademlia) GetFromNetwork(info *ContactInformation, key Key) (interface{}, error) {
	methodName := "GetFromNetwork"
	rpcBase := &util.RPCBase{
		MethodName: methodName,
	}
	args := make([]string, 0)
	args = append(args, contactInfoToString(info))
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

func (kc *RemoteKademlia) Store(info *ContactInformation, key Key, i interface{}) error {
	methodName := "Store"
	rpcBase := &util.RPCBase{
		MethodName: methodName,
	}
	args := make([]interface{}, 0)
	args = append(args, contactInfoToString(info))
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

func (kc *RemoteKademlia) Get(info *ContactInformation, key Key) (interface{}, error) {
	methodName := "Get"
	rpcBase := &util.RPCBase{
		MethodName: methodName,
	}
	args := make([]string, 0)
	args = append(args, contactInfoToString(info))
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

// Get Kademlia nodes from string
func getKademliaNodes(kdNodes string) ([]Kademlia, error) {
	nodes := strings.Split(kdNodes, ".")
	response := make([]Kademlia, 0)
	for _, val := range nodes {
		node := strings.Split(val, ",")
		// node[0] =-> NodeId
		// node[1] =-> Ip
		// node[2] =-> Port
		if len(node) != 3 {
			return nil, errors.New("fail parsing nodes, mismatch params number")
		}
		var key Key
		err := key.GetFromString(node[0])
		if err != nil {
			return nil, errors.New("fail parsing nodes, couldn't get key from string")
		}

		port, err := strconv.Atoi(node[2])
		if err != nil {
			return nil, errors.New("fail parsing nodes, couldn't get port from string")
		}
		response = append(response, NewRemoteKademlia(key, node[1], port))
	}

	return response, nil
}

// Convert Contact Information into string
func contactInfoToString(cInfo *ContactInformation) string {
	return fmt.Sprintf("%v,%v,%v,%v", cInfo.node.GetNodeId().GetString(), cInfo.node.GetIP(), cInfo.node.GetIP(), cInfo.time)
}