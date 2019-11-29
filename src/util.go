package kademlia

func XOR(a uint64, b uint64) uint64 {
	return a ^ b
}

func log(a uint64) int {
	return 1
}

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

type StorageManager interface {
	Store(Key, interface{}) error
	Get(Key) (interface{}, error)
}
