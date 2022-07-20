package dissemination

import (
	"fmt"
	"log"
	"time"

	"github.com/korkmazkadir/ida-gossip/common"
	"github.com/korkmazkadir/ida-gossip/network"
)

func receiveMultipleBlocks(round int, demux *common.Demux, chunkCount int, dataChunkCount int, peerSet *network.PeerSet, leaderCount int, statLogger *common.StatLogger, electedLeaders []int, timeout int) ([]common.Message, []int) {

	chunkChan, err := demux.GetMessageChunkChan(round)
	if err != nil {
		panic(err)
	}

	receiver := newBlockReceiver(leaderCount, chunkCount, dataChunkCount)
	firstChunkReceived := false

	var timeOutChan <-chan time.Time

	// if timeout value is equal or smaller to 0, the chanel will be nil, and
	// all operations will block forever except close operation that panics
	if timeout > 0 {
		timeOutChan = time.After(time.Duration(timeout) * time.Second)
	}

	chunkReceivedMap := make(map[int]bool)
	for _, leader := range electedLeaders {
		chunkReceivedMap[leader] = false
	}

	var extraChunksReceived []common.Chunk

	for !receiver.ReceivedAll() {

		select {
		case c := <-chunkChan:
			{

				if c.Round != round {
					panic(fmt.Errorf("expected round is %d and chunk from round %d", round, c.Round))
				}

				//TODO: considers only one leader
				if c.Issuer != electedLeaders[0] {
					//log.Printf("A chunk is received from previous leader %d, discarding the chunk", c.Issuer)
					extraChunksReceived = append(extraChunksReceived, c)
					continue
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
		case <-timeOutChan:
			// // checks for unresponsive leaders
			// var leadersToRemove []int
			// for leader, value := range chunkReceivedMap {
			// 	if value == false {
			// 		leadersToRemove = append(leadersToRemove, leader)
			// 	}
			// }

			// // this means that at least one leader is unresponsive
			// if len(leadersToRemove) > 0 {
			// 	return nil, leadersToRemove
			// }

			//it enques all unprocessed chunks in case of timeout
			if len(extraChunksReceived) > 0 {
				demux.ReEnqueChunks(extraChunksReceived, round)
				log.Printf("Enqueued %d chunks to reprocess\n", len(extraChunksReceived))
			}

			// WARNING: all leaders are evicted...
			log.Printf("Dissemination Timeout expired: All leaders considered FAULTY!!!")
			return nil, electedLeaders
		}

	}

	for i := range extraChunksReceived {
		log.Printf("Chunk will be discarted: Round: %d Sender: %d\n", extraChunksReceived[i].Round, extraChunksReceived[i].Issuer)
	}

	return receiver.GetBlocks(), nil
}
