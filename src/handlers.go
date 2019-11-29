package kademlia

import (
	"strconv"
)

type HandlerRPC struct {
	kademlia Kademlia
}

func (handler *HandlerRPC) Ping() string {
	return "Pong"
}

func (handler *HandlerRPC) Store(key string, i interface{}) string {
	return ""
}
func (handler *HandlerRPC) Get(key string) string {
	return ""
}
func (handler *HandlerRPC) StoreOnNetwork(key string, i interface{}) string {
	return ""
}
func (handler *HandlerRPC) GetFromNetwork(key string) string {
	return ""
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
func (handler *HandlerRPC) JoinNetwork(key, ip, port string) string {
	return "true"
}
