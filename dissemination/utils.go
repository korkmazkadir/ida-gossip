package dissemination

import (
	"fmt"
	"time"

	"github.com/korkmazkadir/ida-gossip/common"
	"github.com/korkmazkadir/ida-gossip/network"
)

func receiveMultipleBlocks(round int, demux *common.Demux, chunkCount int, dataChunkCount int, peerSet *network.PeerSet, leaderCount int, statLogger *common.StatLogger, electedLeaders []int) ([]common.Message, []int) {

	chunkChan, err := demux.GetMessageChunkChan(round)
	if err != nil {
		panic(err)
	}

	receiver := newBlockReceiver(leaderCount, chunkCount, dataChunkCount)
	firstChunkReceived := false

	//TODO: get timeout value from config
	timeOut := time.After(2 * time.Minute)

	chunkReceivedMap := make(map[int]bool)
	for _, leader := range electedLeaders {
		chunkReceivedMap[leader] = false
	}

	for !receiver.ReceivedAll() {

		select {
		case c := <-chunkChan:
			{

				if c.Round != round {
					panic(fmt.Errorf("expected round is %d and chunk from round %d", round, c.Round))
				}

				// a chunk is received from the leader
				chunkReceivedMap[c.Issuer] = true

				if !firstChunkReceived {
					statLogger.FirstChunkReceived(round, (time.Now().UnixMilli() - c.Time))
					firstChunkReceived = true
				}

				receiver.AddChunk(c)
				peerSet.ForwardChunk(c)
			}
		case <-timeOut:
			// checks for unresponsive leaders
			var leadersToRemove []int
			for leader, value := range chunkReceivedMap {
				if value == false {
					leadersToRemove = append(leadersToRemove, leader)
				}
			}

			// this means that at least one leader is unresponsive
			if len(leadersToRemove) > 0 {
				return nil, leadersToRemove
			}
		}

	}

	return receiver.GetBlocks(), nil
}
