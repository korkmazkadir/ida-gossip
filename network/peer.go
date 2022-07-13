package network

import (
	"errors"

	"github.com/korkmazkadir/ida-gossip/common"
)

var NoCorrectPeerAvailable = errors.New("there are no correct peers available")

type PeerSet struct {
	peers []*P2PClient
}

func (p *PeerSet) AddPeer(IPAddress string, portNumber int, connectionCount int) error {

	client, err := NewClient(IPAddress, portNumber, connectionCount)
	if err != nil {
		return err
	}

	// starts the main loop of client
	go client.Start()

	p.peers = append(p.peers, client)

	return nil
}

func (p *PeerSet) DissaminateChunks(chunks []common.Chunk) {

	for index, chunk := range chunks {
		peer := p.selectPeer(index)
		peer.SendBlockChunk(chunk)
	}
}

func (p *PeerSet) ForwardChunk(chunk common.Chunk) {

	forwardCount := 0
	for _, peer := range p.peers {
		if peer.err != nil {
			continue
		}
		forwardCount++
		peer.SendBlockChunk(chunk)
	}

	if forwardCount == 0 {
		panic(NoCorrectPeerAvailable)
	}
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
