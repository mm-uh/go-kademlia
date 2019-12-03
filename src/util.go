package kademlia

import (
	"errors"
	"fmt"
	"sync"
)

func XOR(a uint64, b uint64) uint64 {
	return a ^ b
}

//func log(a uint64) int {
//	return 1
//}

func Pow(a int, b int) int {
	if b == 0 {
		return 1
	}
	if b == 2 {
		return a * a
	}
	halfPow := Pow(a, b/2)
	if b%2 == 1 {
		return Pow(halfPow, 2) * a
	}
	return Pow(halfPow, 2)
}

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a uint64, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

type StorageManager interface {
	Store(Key, string) error
	Get(Key) (string, error)
}

type ContactInformation struct {
	node Kademlia
	time uint64
}

type NodeInformation struct {
	KBucketNodeInformation
	Kbuckets map[int][]KBucketNodeInformation `json:"KBuckets"`
}

type KBucketNodeInformation struct {
	Key  string `json:"Key"`
	Ip   string `json:"Ip"`
	Port int    `json:"Port"`
}

type SimpleKeyValueStore struct {
	data  map[string]string
	mutex sync.Mutex
}

func NewSimpleKeyValueStore() *SimpleKeyValueStore {
	return &SimpleKeyValueStore{
		data:  make(map[string]string),
		mutex: sync.Mutex{},
	}
}

func (kv *SimpleKeyValueStore) Store(id Key, data string) error {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	kv.data[id.String()] = data
	fmt.Println("SAVING ", data)
	return nil
}

func (kv *SimpleKeyValueStore) Get(id Key) (string, error) {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()
	val, ok := kv.data[id.String()]
	if !ok {
		return "", errors.New("There is no value for that key")
	}
	return val, nil
}
