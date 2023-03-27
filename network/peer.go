package network

import (
	"errors"
	"log"

	"github.com/korkmazkadir/ida-gossip/common"
)

var NoCorrectPeerAvailable = errors.New("there are no correct peers available")

type PeerSet struct {
	peers    []*P2PClient
	fanout   int
	isFaulty bool
	sem      chan int
}

func NewPeerSet(fanout int, parallelSendCount int) *PeerSet {

	var sem chan int
	if parallelSendCount > 0 {
		sem = make(chan int, parallelSendCount)
	}

	p := &PeerSet{
		fanout: fanout,
		sem:    sem,
	}

	log.Printf("PeerSet created. Fanout is %d\n", fanout)
	return p
}

func (p *PeerSet) AddPeer(IPAddress string, portNumber int, connectionCount int) error {

	client, err := NewClient(IPAddress, portNumber, connectionCount, p.sem)
	if err != nil {
		return err
	}

	// starts the main loop of client
	go client.Start()

	p.peers = append(p.peers, client)

	return nil
}

func (p *PeerSet) DissaminateChunks(chunks []common.Chunk) {

	if p.isFaulty {
		return
	}

	for index, chunk := range chunks {
		peer := p.selectPeer(index)
		peer.SendBlockChunk(chunk)
	}
}

func (p *PeerSet) ForwardChunk(chunk common.Chunk) {

	if p.isFaulty {
		return
	}

	// Provides a random sampling for each chunk with in 16 peers.
	//rand.Seed(time.Now().UnixNano())
	//rand.Shuffle(len(p.peers), func(i, j int) { p.peers[i], p.peers[j] = p.peers[j], p.peers[i] })

	forwardCount := 0
	for _, peer := range p.peers {
		if peer.err != nil {
			continue
		}

		peer.SendBlockChunk(chunk)
		forwardCount++

		if forwardCount == p.fanout {
			// the message is forwarded fanout times so break here
			break
		}

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

func (p *PeerSet) SetFaulty() {
	p.isFaulty = true
}
