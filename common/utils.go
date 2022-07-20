package common

import (
	"crypto/rand"
	"encoding/base64"
	"math"
)

func EncodeBase64(hex []byte) string {
	return base64.StdEncoding.EncodeToString([]byte(hex))
}

func GetRandomByteSlice(size int) []byte {
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	return data
}

func FaultyNodeCount(nodeCount int, faultyNodePercent int) int {

	if nodeCount <= 0 {
		panic("node count is smaller or equal to 0")
	}

	if faultyNodePercent <= 0 {
		return 0
	}

	return int(math.Floor((float64(nodeCount) / float64(100)) * float64(faultyNodePercent)))
}

func IsFaulty(nodeCount int, faultyNodePercent int, nodeID int) bool {

	if nodeCount <= 0 {
		panic("node count is smaller or equal to 0")
	}

	if faultyNodePercent <= 0 {
		return false
	}

	return nodeID >= (nodeCount - FaultyNodeCount(nodeCount, faultyNodePercent))
}
