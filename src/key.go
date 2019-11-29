package kademlia

import "errors"

type KeyNode256 [32]byte

func (kn *KeyNode256) XOR(other Key) (Key, error) {
	otherKeyNode, ok := other.(*KeyNode256)
	if !ok {
		return nil, errors.New("Other is not a valid type")
	}

	var result KeyNode

	for i := 0; i < 32; i++ {
		result[i] = kn[i] ^ otherKeyNode[i]
	}
	return &result, nil
}

func (kn *KeyNode256) IsActive(index int) bool {
	nByte := index / 8
	nBit := index % 8

	return (kn[nByte] & uint8(Pow(2, nBit))) != 0

}

func (kn *KeyNode256) Lenght() int {
	return 256
}

func (kn *KeyNode256) Less(other interface{}) (bool, error) {
	otherKeyNode, ok := other.(*KeyNode256)
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

type KeyNode [20]byte

func (kn *KeyNode) XOR(other Key) (Key, error) {
	otherKeyNode, ok := other.(*KeyNode)
	if !ok {
		return nil, errors.New("Other is not a valid type")
	}

	var result KeyNode

	for i := 0; i < 20; i++ {
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
	return 160
}

func (kn *KeyNode) Less(other interface{}) (bool, error) {
	otherKeyNode, ok := other.(*KeyNode)
	if !ok {
		return false, errors.New("Other is not a valid type")
	}

	for i := 19; i >= 0; i-- {
		if kn[i] == otherKeyNode[i] {
			continue
		}
		return kn[i] < otherKeyNode[i], nil
	}
	return false, nil
}
