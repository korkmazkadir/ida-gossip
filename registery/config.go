package registery

import (
	"crypto/sha256"
	"fmt"
)

type NodeConfig struct {
	NodeCount int

	EpochSeed []byte

	EndRound int

	RoundSleepTime int

	GossipFanout int

	ConnectionCount int

	SourceCount int

	MessageSize int

	MessageChunkCount int
}

func (nc NodeConfig) Hash() []byte {

	str := fmt.Sprintf("%d,%x,%d,%d,%d,%d,%d,%d", nc.NodeCount, nc.EpochSeed, nc.EndRound, nc.GossipFanout, nc.SourceCount, nc.MessageSize, nc.MessageChunkCount, nc.ConnectionCount)

	h := sha256.New()
	_, err := h.Write([]byte(str))
	if err != nil {
		panic(err)
	}

	return h.Sum(nil)
}

func (nc *NodeConfig) CopyFields(cp NodeConfig) {
	nc.NodeCount = cp.NodeCount
	nc.EpochSeed = nc.EpochSeed[:0]
	nc.EpochSeed = append(nc.EpochSeed, cp.EpochSeed...)
	nc.EndRound = cp.EndRound
	nc.RoundSleepTime = cp.RoundSleepTime
	nc.GossipFanout = cp.GossipFanout
	nc.ConnectionCount = cp.ConnectionCount
	nc.SourceCount = cp.SourceCount
	nc.MessageSize = cp.MessageSize
	nc.MessageChunkCount = cp.MessageChunkCount
}
