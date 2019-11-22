package kademlia

import "errors"

type KeyNode [32]byte

func (kn *KeyNode) XOR(other Key) (Key, error) {
	otherKeyNode, ok := other.(*KeyNode)
	if !ok {
		return nil, errors.New("Other is not a valid type")
	}

	var result KeyNode

	for i := 0; i < 32; i++ {
		result[i] = kn[i] ^ otherKeyNode[i]
	}
	return &result, nil
}

func (kn *KeyNode) IsActive(index int) bool {
	nByte := index / 8
	nBit := index % 8

	return (kn[nByte] & uint8(Pow(2, nBit))) != 0

}

func (kn *KeyNode) Lenght() int {
	return 256
}

func (kn *KeyNode) Less(other Key) (bool, error) {
	otherKeyNode, ok := other.(*KeyNode)
	if !ok {
		return false, errors.New("Other is not a valid type")
	}

	for i := 31; i >= 0; i-- {
		if kn[i] == otherKeyNode[i] {
			continue
		}
		return kn[i] < otherKeyNode[i], nil
	}
	return false, nil
}

//ToDo
func KeyNodeFromSHA256() *KeyNode {
	return nil
}
