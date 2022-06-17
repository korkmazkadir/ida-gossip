package dissemination

import (
	"github.com/korkmazkadir/ida-gossip/common"
	"github.com/korkmazkadir/ida-gossip/network"
)

func receiveMultipleBlocks(round int, demux *common.Demux, chunkCount int, peerSet *network.PeerSet, leaderCount int) []common.Message {

	chunkChan, err := demux.GetVoteBlockChunkChan(round)
	if err != nil {
		panic(err)
	}

	receiver := newBlockReceiver(leaderCount, chunkCount)
	for !receiver.ReceivedAll() {
		c := <-chunkChan
		receiver.AddChunk(c)
		peerSet.ForwardChunk(c)
	}

	return receiver.GetBlocks()
}
