package kademlia

import (
	"strconv"
)

type HandlerRPC struct {
	kademlia Kademlia
}

func (handler *HandlerRPC) Ping() string {
	return strconv.FormatBool(handler.kademlia.Ping())
}

func (handler *HandlerRPC) Store(keyAsString string, i interface{}) string {
	var key Key
	err := key.GetFromString(keyAsString)
	if err != nil {
		return "false"
	}
	err = handler.kademlia.Store(key, i)
	if err != nil {
		return "false"
	}
	return "true"
}

func (handler *HandlerRPC) Get(keyAsString string) string {
	var key Key
	err := key.GetFromString(keyAsString)
	if err != nil {
		return "false"
	}
	val, err := handler.kademlia.Get(key)
	if err != nil {
		return "false"
	}
	str, ok := val.(string)
	if !ok {
		return "false"
	}
	return str
}

func (handler *HandlerRPC) StoreOnNetwork(keyAsString string, i interface{}) string {
	var key Key
	err := key.GetFromString(keyAsString)
	if err != nil {
		return "false"
	}
	err = handler.kademlia.StoreOnNetwork(key, i)
	if err != nil {
		return "false"
	}
	return "true"
}

func (handler *HandlerRPC) GetFromNetwork(keyAsString string) string {
	var key Key
	err := key.GetFromString(keyAsString)
	if err != nil {
		return "false"
	}
	val, err := handler.kademlia.GetFromNetwork(key)
	if err != nil {
		return "false"
	}
	str, ok := val.(string)
	if !ok {
		return "false"
	}
	return str
}
func (handler *HandlerRPC) GetNodeId() string {
	return handler.kademlia.GetNodeId().GetString()
}
func (handler *HandlerRPC) GetIP() string {
	return handler.kademlia.GetIP()
}
func (handler *HandlerRPC) GetPort() string {
	return strconv.FormatInt(int64(handler.kademlia.GetPort()), 10)
}
func (handler *HandlerRPC) JoinNetwork(keyAsString, ip, portAsString string) string {
	var key Key
	err := key.GetFromString(keyAsString)
	if err != nil {
		return "false"
	}
	port, err := strconv.Atoi(portAsString)
	if err != nil {
		return "false"
	}
	kdToJoin := NewKademliaContact(key, ip, port)
	err = handler.kademlia.JoinNetwork(kdToJoin)
	if err != nil {
		return "false"
	}
	return "true"
}
