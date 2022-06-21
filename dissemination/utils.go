package dissemination

import (
	"fmt"
	"time"

	"github.com/korkmazkadir/ida-gossip/common"
	"github.com/korkmazkadir/ida-gossip/network"
)

func receiveMultipleBlocks(round int, demux *common.Demux, chunkCount int, peerSet *network.PeerSet, leaderCount int, statLogger *common.StatLogger) []common.Message {

	chunkChan, err := demux.GetMessageChunkChan(round)
	if err != nil {
		panic(err)
	}

	receiver := newBlockReceiver(leaderCount, chunkCount)
	firstChunkReceived := false
	for !receiver.ReceivedAll() {
		c := <-chunkChan

		if c.Round != round {
			panic(fmt.Errorf("expected round is %d and chunk from round %d", round, c.Round))
		}

		if !firstChunkReceived {
			statLogger.FirstChunkReceived(round, (time.Now().UnixMilli() - c.Time))
			firstChunkReceived = true
		}

		receiver.AddChunk(c)
		peerSet.ForwardChunk(c)
	}

	return receiver.GetBlocks()
}
