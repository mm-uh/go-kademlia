package kademlia

import (
	"strconv"
)

type HandlerRPC struct {
	kademlia Kademlia
}

func (handler *HandlerRPC) Ping(cInfo string) string {
	return strconv.FormatBool(handler.kademlia.Ping(getContactInformationFromString(cInfo)))
}

func (handler *HandlerRPC) Store(cInfo string, keyAsString string, i interface{}) string {
	var key Key
	err := key.GetFromString(keyAsString)
	if err != nil {
		return "false"
	}
	err = handler.kademlia.Store(getContactInformationFromString(cInfo), key, i)
	if err != nil {
		return "false"
	}
	return "true"
}

func (handler *HandlerRPC) Get(cInfo, keyAsString string) string {
	var key Key
	err := key.GetFromString(keyAsString)
	if err != nil {
		return "false"
	}
	val, err := handler.kademlia.Get(getContactInformationFromString(cInfo), key)
	if err != nil {
		return "false"
	}
	str, ok := val.(string)
	if !ok {
		return "false"
	}
	return str
}

func (handler *HandlerRPC) StoreOnNetwork(cInfo, keyAsString string, i interface{}) string {
	var key Key
	err := key.GetFromString(keyAsString)
	if err != nil {
		return "false"
	}
	err = handler.kademlia.StoreOnNetwork(getContactInformationFromString(cInfo), key, i)
	if err != nil {
		return "false"
	}
	return "true"
}

func (handler *HandlerRPC) GetFromNetwork(cInfo, keyAsString string) string {
	var key Key
	err := key.GetFromString(keyAsString)
	if err != nil {
		return "false"
	}
	val, err := handler.kademlia.GetFromNetwork(getContactInformationFromString(cInfo), key)
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
	kdToJoin := NewRemoteKademlia(key, ip, port)
	err = handler.kademlia.JoinNetwork(kdToJoin)
	if err != nil {
		return "false"
	}
	return "true"
}

func getContactInformationFromString(string) *ContactInformation {
	// TODO Update this method to convert Contact Information String into proper Struct
	return &ContactInformation{}
}
