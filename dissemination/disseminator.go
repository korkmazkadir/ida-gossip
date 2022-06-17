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

	// starts a new epoch
	d.statLogger.NewRound(round)

	// sets the round for demultiplexer
	d.demultiplexer.UpdateRound(round)

	// chunks the block
	chunks := common.ChunkMessage(message, d.nodeConfig.BlockChunkCount)
	//log.Printf("proposing block %x\n", encodeBase64(merkleRoot[:15]))

	// disseminate chunks over different nodes
	d.peerSet.DissaminateChunks(chunks)

	//return d.WaitForMessage(round)
}

func (d *Disseminator) WaitForMessage(round int) []common.Message {

	// starts a new epoch
	d.statLogger.NewRound(round)

	// sets the round for demultiplexer
	d.demultiplexer.UpdateRound(round)

	startTime := time.Now()
	messages := receiveMultipleBlocks(round, d.demultiplexer, d.nodeConfig.BlockChunkCount, &d.peerSet, d.nodeConfig.LeaderCount)
	d.statLogger.LogBlockReceive(time.Since(startTime).Milliseconds())

	d.statLogger.LogEndOfRound()

	return messages
}
