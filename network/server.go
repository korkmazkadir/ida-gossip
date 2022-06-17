package network

import (
	"github.com/korkmazkadir/ida-gossip/common"
)

type P2PServer struct {
	demux *common.Demux
}

func NewServer(demux *common.Demux) *P2PServer {
	server := &P2PServer{demux: demux}
	return server
}

func (s *P2PServer) HandleBlockChunk(chunk *common.Chunk, reply *int) error {

	s.demux.EnqueBlockChunk(*chunk)

	return nil
}
