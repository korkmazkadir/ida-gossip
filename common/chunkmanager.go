package common

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
)

func ChunkMessage(message Message, numberOfChunks int) []Chunk {

	blockBytes := encodeToBytes(message)
	chunks := constructChunks(message, blockBytes, numberOfChunks)
	return chunks
}

// mergeChunks assumes that sanity checks are done before calling this function
func MergeChunks(chunks []Chunk) Message {

	var blockData []byte
	for i := 0; i < len(chunks); i++ {
		blockData = append(blockData, chunks[i].Payload...)
	}

	return decodeToBlock(blockData)
}

func constructChunks(message Message, blockBytes []byte, numberOfChunks int) []Chunk {

	var chunks []Chunk
	chunkSize := int(math.Ceil(float64(len(blockBytes)) / float64(numberOfChunks)))

	if chunkSize == 0 {
		panic(fmt.Errorf("chunk payload size is 0"))
	}

	for i := 0; i < numberOfChunks; i++ {

		startIndex := i * chunkSize
		endIndex := startIndex + chunkSize

		var payload []byte
		if i < (numberOfChunks - 1) {
			payload = blockBytes[startIndex:endIndex]
		} else {
			payload = blockBytes[startIndex:]
		}

		chunk := Chunk{
			Round:      message.Round,
			Time:       message.Time,
			ChunkCount: numberOfChunks,
			ChunkIndex: i,
			Payload:    payload,
		}

		chunks = append(chunks, chunk)
	}

	return chunks
}

// https://gist.github.com/SteveBate/042960baa7a4795c3565
func encodeToBytes(p interface{}) []byte {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func decodeToBlock(data []byte) Message {

	message := Message{}
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&message)
	if err != nil {
		panic(err)
	}
	return message
}
