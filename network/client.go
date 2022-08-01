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

	//rpcClient *rpc.Client

	rpcClients []*rpc.Client
	err        error
}

// NewClient creates a new client
func NewClient(IPAddress string, portNumber int, connectionCount int) (*P2PClient, error) {

	if connectionCount < 1 {
		panic(fmt.Errorf("connection count is %d, it must be bigger than 1", connectionCount))
	}

	var clients []*rpc.Client
	for i := 0; i < connectionCount; i++ {
		rpcClient, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", IPAddress, portNumber))
		if err != nil {
			return nil, err
		}
		clients = append(clients, rpcClient)
	}

	client := &P2PClient{}
	client.IPAddress = IPAddress
	client.portNumber = portNumber
	client.rpcClients = append(client.rpcClients, clients...)

	return client, nil
}

// SendBlockChunk enqueues a chunk of a block to send
func (c *P2PClient) SendBlockChunk(chunk common.Chunk) {

	err := c.rpcClients[0].Call("P2PServer.HandleBlockChunk", chunk, nil)
	//TODO: needs to handle the error properly
	if err != nil {
		c.err = err
		panic(err)
	}
}
