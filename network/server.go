package network

import (
	"github.com/korkmazkadir/ida-gossip/common"
)

type P2PServer struct {
	demux    *common.Demux
	isFaulty bool
}

func NewServer(demux *common.Demux) *P2PServer {
	server := &P2PServer{demux: demux, isFaulty: false}
	return server
}

func (s *P2PServer) SetAsFaulty() {
	s.isFaulty = true
}

func (s *P2PServer) HandleBlockChunk(chunk *common.Chunk, reply *int) error {

	if s.isFaulty {
		// if it is faulty ignores all messages
		return nil
	}

	s.demux.EnqueBlockChunk(*chunk)

	return nil
}
