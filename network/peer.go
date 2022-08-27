package network

import (
	"errors"
	"math"
	"sync"
	"time"

	"github.com/korkmazkadir/ida-gossip/common"
)

var NoCorrectPeerAvailable = errors.New("there are no correct peers available")

type sendStats struct {
	mutex sync.Mutex
	stats []int64
}

func (s *sendStats) append(stat int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.stats = append(s.stats, stat)
}

func (s *sendStats) summary() (min int64, mean float64, max int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	min = math.MaxInt64

	var sum int64
	for _, value := range s.stats {

		sum += value
		if value < min {
			min = value
		}

		if value > max {
			max = value
		}
	}

	mean = float64(sum) / float64(len(s.stats))

	return
}

type PeerSet struct {
	peers               []*P2PClient
	jobChan             chan job
	isFaulty            bool
	concurrentSendCount int
	stats               *sendStats
}

func NewPeerSet(concurrentSendCount int) *PeerSet {
	p := &PeerSet{
		jobChan:             make(chan job, 1024),
		concurrentSendCount: concurrentSendCount,
		stats:               &sendStats{},
	}
	return p
}

func (p *PeerSet) AddPeer(IPAddress string, portNumber int, connectionCount int) error {

	client, err := NewClient(IPAddress, portNumber, connectionCount)
	if err != nil {
		return err
	}

	p.peers = append(p.peers, client)

	return nil
}

func (p *PeerSet) DissaminateChunks(chunks []common.Chunk) {

	if p.isFaulty {
		return
	}

	for index, chunk := range chunks {
		peer := p.selectPeer(index)
		j := job{chunk: chunk, peer: peer}
		p.jobChan <- j
	}
}

func (p *PeerSet) ForwardChunk(chunk common.Chunk) {

	if p.isFaulty {
		return
	}

	j := job{chunk: chunk}
	p.jobChan <- j
}

func (p *PeerSet) selectPeer(index int) *P2PClient {

	peerCount := len(p.peers)
	for i := 0; i < peerCount; i++ {
		peer := p.peers[(index+i)%peerCount]
		if peer.err == nil {
			return peer
		}
	}

	panic(NoCorrectPeerAvailable)
}

func (p *PeerSet) SetFaulty() {
	p.isFaulty = true
}

type job struct {
	chunk common.Chunk
	peer  *P2PClient
}

func (p *PeerSet) MainLoop() {

	var sem = make(chan int, p.concurrentSendCount)

	for {

		j := <-p.jobChan

		if j.peer != nil {

			sem <- 1
			go func() {
				start := time.Now()
				j.peer.SendBlockChunk(j.chunk)
				<-sem
				p.stats.append(time.Since(start).Milliseconds())
			}()

			continue
		}

		for _, peer := range p.peers {
			// if it can not send, it panics
			// no need to check for errors

			sem <- 1
			go func(peer *P2PClient) {
				start := time.Now()
				peer.SendBlockChunk(j.chunk)
				<-sem
				p.stats.append(time.Since(start).Milliseconds())
			}(peer)

		}

	}

}

func (p *PeerSet) SendStats() (int64, float64, int64) {

	return p.stats.summary()
}
