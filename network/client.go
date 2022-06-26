package network

import (
	"fmt"
	"net/rpc"
	"sync"
	"sync/atomic"

	"github.com/korkmazkadir/ida-gossip/common"
)

// Client implements P2P client
type P2PClient struct {
	IPAddress  string
	portNumber int

	rpcClient *rpc.Client

	blockChunks chan common.Chunk

	mutexQueueLength  sync.Mutex
	sumQueueLength    int32
	queueLengthCount  int32
	onAirMessageCount int32

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

// SendBlockChunk enqueues a chunk of a block to send
func (c *P2PClient) SendBlockChunk(chunk common.Chunk) {

	c.blockChunks <- chunk

	c.mutexQueueLength.Lock()
	defer c.mutexQueueLength.Unlock()
	//calculates average queue length
	c.sumQueueLength += atomic.LoadInt32(&c.onAirMessageCount)
	c.queueLengthCount++
}

func (c *P2PClient) ResetQueueLengthCounters() {

	c.mutexQueueLength.Lock()
	defer c.mutexQueueLength.Unlock()

	c.queueLengthCount = 0
	c.sumQueueLength = 0
}

func (c *P2PClient) GetAvgQueueLength() float64 {

	c.mutexQueueLength.Lock()
	defer c.mutexQueueLength.Unlock()

	return float64(c.sumQueueLength) / float64(c.queueLengthCount)
}

func (c *P2PClient) mainLoop() {

	for {
		blockChunk := <-c.blockChunks
		//go c.rpcClient.Call("P2PServer.HandleBlockChunk", blockChunk, nil)

		go func() {

			atomic.AddInt32(&c.onAirMessageCount, 1)

			c.rpcClient.Call("P2PServer.HandleBlockChunk", blockChunk, nil)

			atomic.AddInt32(&c.onAirMessageCount, -1)

		}()

	}
}
