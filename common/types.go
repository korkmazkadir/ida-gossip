package common

import (
	"crypto/sha256"
	"fmt"
)

// Block defines blockchain block structure
type Message struct {
	Issuer int

	Round int

	Time int64

	Payload []byte
}

func (m Message) Hash() []byte {

	h := sha256.New()
	//writing paylaod first seems more memory efficient because no need to create a long string
	_, err := h.Write(m.Payload)
	if err != nil {
		panic(err)
	}

	str := fmt.Sprintf("%d,%d", m.Issuer, m.Round)
	_, err = h.Write([]byte(str))
	if err != nil {
		panic(err)
	}

	return h.Sum(nil)
}

// BlockChunk defines a chunk of a block.
// BlockChunks disseminate fater in the gossip network because they are very small compared to a Block
type Chunk struct {
	Issuer int

	// Round of the block
	Round int

	Time int64

	// The number of expected chunks to reconstruct a block
	ChunkCount int

	// Chunk index
	ChunkIndex int

	// Chunk payload
	Payload []byte
}

// Hash produces the digest of a BlockChunk.
// It considers all fields of a BlockChunk.
func (c Chunk) Hash() []byte {

	h := sha256.New()
	//writing paylaod first seems more memory efficient because no need to create a long string
	_, err := h.Write(c.Payload)
	if err != nil {
		panic(err)
	}

	str := fmt.Sprintf("%d,%d,%d,%d", c.Round, c.ChunkCount, c.ChunkIndex, c.Issuer)
	_, err = h.Write([]byte(str))
	if err != nil {
		panic(err)
	}

	return h.Sum(nil)
}
