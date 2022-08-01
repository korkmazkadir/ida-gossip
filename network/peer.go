package network

import (
	"errors"

	"github.com/korkmazkadir/ida-gossip/common"
)

var NoCorrectPeerAvailable = errors.New("there are no correct peers available")

type PeerSet struct {
	peers    []*P2PClient
	jobChan  chan job
	isFaulty bool
}

func NewPeerSet() *PeerSet {
	p := &PeerSet{jobChan: make(chan job, 1024)}
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

	for {

		j := <-p.jobChan

		if j.peer != nil {
			j.peer.SendBlockChunk(j.chunk)
			continue
		}

		for _, peer := range p.peers {
			// if it can not send, it panics
			// no need to check for errors
			peer.SendBlockChunk(j.chunk)
		}

	}

}
