package dissemination

import (
	"fmt"
	"sort"

	"github.com/korkmazkadir/ida-gossip/common"
)

type blockReceiver struct {
	blockCount         int
	chunkCount         int
	dataChunkCount     int
	blockMap           map[string][]common.Chunk
	receivedChunkCount int
}

func newBlockReceiver(leaderCount int, chunkCount int, dataChunkCount int) *blockReceiver {

	r := &blockReceiver{
		blockCount:     leaderCount,
		chunkCount:     chunkCount,
		dataChunkCount: dataChunkCount,
		blockMap:       make(map[string][]common.Chunk),
	}

	return r
}

// AddChunk stores a chunk of a block to reconstruct the whole block later
func (r *blockReceiver) AddChunk(chunk common.Chunk) {
	key := string(chunk.Issuer)
	chunkSlice := r.blockMap[key]
	r.blockMap[key] = append(chunkSlice, chunk)

	//r.receivedChunkCount++
	//log.Printf("A chunk received: sender %d index %d total count %d\n", chunk.Issuer, chunk.ChunkIndex, r.receivedChunkCount)

	if len(r.blockMap) > r.blockCount {
		panic(fmt.Errorf("there are more blocks than expected, the number of blocks is %d", len(r.blockMap)))
	}
}

// ReceivedAll checks whether all chunks are recived or not to reconstruct the blocks of a round
func (r *blockReceiver) ReceivedAll() bool {

	if len(r.blockMap) != r.blockCount {
		return false
	}

	for _, chunkSlice := range r.blockMap {
		if len(chunkSlice) != r.dataChunkCount {
			return false
		}
	}

	return true
}

// GetBlocks recunstruct blocks using chunks, and returns the blocks by sorting the resulting block slice according to block hashes
func (r *blockReceiver) GetBlocks() []common.Message {

	if r.ReceivedAll() == false {
		panic(fmt.Errorf("not received all block chunks to reconstruct block/s"))
	}

	keys := make([]string, 0, len(r.blockMap))
	for k := range r.blockMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var messages []common.Message

	for _, key := range keys {

		receivedChunks := r.blockMap[key]
		sort.Slice(receivedChunks, func(i, j int) bool {
			return receivedChunks[i].ChunkIndex < receivedChunks[j].ChunkIndex
		})

		block := common.MergeChunks(receivedChunks, r.chunkCount, r.dataChunkCount)
		messages = append(messages, block)
	}

	return messages
}
