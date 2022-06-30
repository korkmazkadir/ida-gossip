package dissemination

import (
	"errors"
	"time"

	"github.com/korkmazkadir/ida-gossip/common"
	"github.com/korkmazkadir/ida-gossip/network"
	"github.com/korkmazkadir/ida-gossip/registery"
)

// BlockNotValid is returned if the block can not pass vaslidity test
var ErrBlockNotValid = errors.New("received block is not valid")

type Disseminator struct {
	demultiplexer *common.Demux
	nodeConfig    registery.NodeConfig
	peerSet       network.PeerSet
	statLogger    *common.StatLogger
}

func NewDisseminator(demux *common.Demux, config registery.NodeConfig, peerSet network.PeerSet, statLogger *common.StatLogger) *Disseminator {

	d := &Disseminator{
		demultiplexer: demux,
		nodeConfig:    config,
		peerSet:       peerSet,
		statLogger:    statLogger,
	}

	return d
}

func (d *Disseminator) SubmitMessage(round int, message common.Message) {

	// sets the round for demultiplexer
	d.demultiplexer.UpdateRound(round)

	// chunks the block
	chunks := common.ChunkMessage(message, d.nodeConfig.MessageChunkCount)
	//log.Printf("proposing block %x\n", encodeBase64(merkleRoot[:15]))

	// disseminate chunks over different nodes
	d.peerSet.DissaminateChunks(chunks)

	//return d.WaitForMessage(round)
}

func (d *Disseminator) WaitForMessage(round int, electedLeaders []int) ([]common.Message, []int) {

	// sets the round for demultiplexer
	d.demultiplexer.UpdateRound(round)

	messages, leadersToRemove := receiveMultipleBlocks(round, d.demultiplexer, d.nodeConfig.MessageChunkCount, &d.peerSet, d.nodeConfig.SourceCount, d.statLogger, electedLeaders)

	if leadersToRemove != nil {
		d.peerSet.ResetQueueLengthCounters()
		return nil, leadersToRemove
	}

	d.statLogger.AvgQueuLength(round, d.peerSet.GetAvgQueueLength())
	d.peerSet.ResetQueueLengthCounters()

	var maxElapsedTime int64
	for i := range messages {
		elapsedTime := time.Now().UnixMilli() - messages[i].Time
		if elapsedTime > maxElapsedTime {
			maxElapsedTime = elapsedTime
		}
	}
	d.statLogger.MessageReceived(round, (maxElapsedTime))

	return messages, nil
}
