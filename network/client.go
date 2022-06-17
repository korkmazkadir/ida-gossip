package network

import (
	"fmt"
	"net/rpc"

	"github.com/korkmazkadir/ida-gossip/common"
)

// Client implements P2P client
type P2PClient struct {
	IPAddress  string
	portNumber int

	rpcClient *rpc.Client

	blockChunks chan common.Chunk

	err error
}

// NewClient creates a new client
func NewClient(IPAddress string, portNumber int) (*P2PClient, error) {

	rpcClient, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", IPAddress, portNumber))
	if err != nil {
		return nil, err
	}

	client := &P2PClient{}
	client.IPAddress = IPAddress
	client.portNumber = portNumber
	client.rpcClient = rpcClient

	client.blockChunks = make(chan common.Chunk, 1024)

	return client, nil
}

// Start starts the main loop of client. It blocks the calling goroutine
func (c *P2PClient) Start() {

	c.mainLoop()
}

// SendBlockChunk enques a chunk of a block to send
func (c *P2PClient) SendBlockChunk(chunk common.Chunk) {

	c.blockChunks <- chunk
}

func (c *P2PClient) mainLoop() {

	for {
		blockChunk := <-c.blockChunks
		go c.rpcClient.Call("P2PServer.HandleBlockChunk", blockChunk, nil)
	}
}
