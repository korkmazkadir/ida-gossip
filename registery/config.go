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

	ParallelSendCount int

	SourceCount int

	MessageSize int

	MessageChunkCount int

	DataChunkCount int

	EndOfExperimentSleepTime int

	FaultyNodePercent int

	DisseminationTimeout int
}

func (nc NodeConfig) Hash() []byte {

	str := fmt.Sprintf("%d,%x,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d", nc.NodeCount, nc.EpochSeed, nc.EndRound, nc.GossipFanout, nc.SourceCount, nc.MessageSize, nc.MessageChunkCount, nc.ConnectionCount, nc.DataChunkCount, nc.EndOfExperimentSleepTime, nc.FaultyNodePercent, nc.DisseminationTimeout, nc.ParallelSendCount)

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
	nc.ParallelSendCount = cp.ParallelSendCount
	nc.SourceCount = cp.SourceCount
	nc.MessageSize = cp.MessageSize
	nc.MessageChunkCount = cp.MessageChunkCount
	nc.DataChunkCount = cp.DataChunkCount
	nc.EndOfExperimentSleepTime = cp.EndOfExperimentSleepTime
	nc.FaultyNodePercent = cp.FaultyNodePercent
	nc.DisseminationTimeout = cp.DisseminationTimeout
}
