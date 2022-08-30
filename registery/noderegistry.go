package registery

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/korkmazkadir/ida-gossip/common"
)

// NodeInfo keeps node info
type NodeInfo struct {
	ID         int
	IPAddress  string
	PortNumber int
}

type NodeList struct {
	Nodes []NodeInfo
}

type NodeRegistry struct {
	mutex           sync.Mutex
	registeredNodes []NodeInfo

	failedNodes   []NodeInfo
	finishedNodes []NodeInfo
	startedNodes  []NodeInfo

	config       NodeConfig
	statKeeper   *StatKeeper
	statusLogger *StatusLogger
	nodeIDs      []int

	once sync.Once
}

func NewNodeRegistry(config NodeConfig, statusLogger *StatusLogger) *NodeRegistry {

	var nodeIDList []int
	for i := 1; i <= config.NodeCount; i++ {
		nodeIDList = append(nodeIDList, i)
	}
	// Shuffles node id list
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(nodeIDList), func(i, j int) {
		nodeIDList[i], nodeIDList[j] = nodeIDList[j], nodeIDList[i]
	})

	return &NodeRegistry{config: config, statusLogger: statusLogger, nodeIDs: nodeIDList}
}

// Register registers a node with specific node info
func (nr *NodeRegistry) Register(nodeInfo *NodeInfo, reply *NodeInfo) error {

	nr.mutex.Lock()
	defer nr.mutex.Unlock()

	// assigns a node ID. smallest node ID is 1
	// nodeID := len(nr.registeredNodes) + 1

	// assigns a random node ID. smallest node ID is 1
	nodeID := nr.nodeIDs[len(nr.registeredNodes)]
	nodeInfo.ID = nodeID

	nr.registeredNodes = append(nr.registeredNodes, *nodeInfo)
	log.Printf("new node registered; ip address %s port number %d, registered node count: %d\n", nodeInfo.IPAddress, nodeInfo.PortNumber, len(nr.registeredNodes))

	reply.IPAddress = nodeInfo.IPAddress
	reply.PortNumber = nodeInfo.PortNumber
	reply.ID = nodeInfo.ID

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////
func (nr *NodeRegistry) NodeStarted(nodeInfo *NodeInfo, reply *int) error {

	nr.mutex.Lock()
	defer nr.mutex.Unlock()

	log.Printf("node started %d\n", nodeInfo.ID)

	nr.startedNodes = append(nr.startedNodes, *nodeInfo)

	startedNodeCount := len(nr.startedNodes)
	failedNodeCount := len(nr.failedNodes)

	if (startedNodeCount + failedNodeCount) == nr.config.NodeCount {
		nr.statusLogger.LogStarted()
	} else {
		log.Printf("%d failed %d started\n", failedNodeCount, startedNodeCount)
	}

	return nil
}

func (nr *NodeRegistry) NodeFailed(nodeInfo *NodeInfo, reply *int) error {

	nr.mutex.Lock()
	defer nr.mutex.Unlock()

	log.Printf("node failed %d\n", nodeInfo.ID)

	nr.failedNodes = append(nr.failedNodes, *nodeInfo)

	failedNodeCount := len(nr.failedNodes)
	faultyNodeCount := common.FaultyNodeCount(nr.config.NodeCount, nr.config.FaultyNodePercent)

	if failedNodeCount-faultyNodeCount >= 20 {
		nr.statusLogger.LogFailed()
		panic(fmt.Errorf("more than 20 nodes failed there must be a problem"))
	}

	return nil
}

func (nr *NodeRegistry) NodeFinished(nodeInfo *NodeInfo, reply *int) error {

	nr.mutex.Lock()
	defer nr.mutex.Unlock()

	log.Printf("node finished %d\n", nodeInfo.ID)

	// auto close in 60 seconds
	//if nr.config.FaultyNodePercent > 0 {
	// TODO: The behavior is changed, be careful!
	nr.once.Do(func() {
		go nr.CountDownToClose()
	})
	//}

	nr.finishedNodes = append(nr.finishedNodes, *nodeInfo)

	finishedNodeCount := len(nr.finishedNodes)
	failedNodeCount := len(nr.failedNodes)

	if (finishedNodeCount + failedNodeCount) == nr.config.NodeCount {
		nr.statusLogger.LogFinished()
	} else {
		log.Printf("%d failed %d finished\n", failedNodeCount, finishedNodeCount)
	}

	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////

func (nr *NodeRegistry) Unregister(remoteAddress string) {
	addressParts := strings.Split(remoteAddress, ":")

	if len(addressParts) != 2 {
		log.Printf("unknown address format, node couldnot unregistered %s \n", remoteAddress)
		return
	}

	ipAddress := addressParts[0]
	portNumber, err := strconv.Atoi(addressParts[1])
	if err != nil {
		log.Printf("could not parse the port number, error: %s, portnumber: %s\n", err, addressParts[1])
		return
	}

	nr.mutex.Lock()
	defer nr.mutex.Unlock()

	nodeIndex := -1
	for i := range nr.registeredNodes {
		if nr.registeredNodes[i].IPAddress == ipAddress && nr.registeredNodes[i].PortNumber == portNumber {
			nodeIndex = i
			break
		}
	}

	if nodeIndex == -1 {
		log.Printf("could not find %s in the registered node list to unregister\n", remoteAddress)
		return
	}

	nr.registeredNodes = append(nr.registeredNodes[:nodeIndex], nr.registeredNodes[nodeIndex+1:]...)
	log.Printf("node %s unregistered successfully\n", remoteAddress)

}

// GetConfig is used to get config
func (nr *NodeRegistry) GetConfig(nodeInfo *NodeInfo, config *NodeConfig) error {

	nr.mutex.Lock()
	defer nr.mutex.Unlock()

	config.CopyFields(nr.config)

	return nil
}

// GetNodeList returns node list
func (nr *NodeRegistry) GetNodeList(nodeInfo *NodeInfo, nodeList *NodeList) error {

	nr.mutex.Lock()
	defer nr.mutex.Unlock()

	nodeList.Nodes = append(nodeList.Nodes, nr.registeredNodes...)

	return nil
}

func (nr *NodeRegistry) UploadStats(stats *common.StatList, reply *int) error {

	nr.mutex.Lock()
	defer nr.mutex.Unlock()

	log.Printf("node %d (%s:%d) uploading stats; event count %d \n", stats.NodeID, stats.IPAddress, stats.PortNumber, len(stats.Events))

	if nr.statKeeper == nil {
		nr.statKeeper = NewStatKeeper(nr.config)
	}

	nr.statKeeper.SaveStats(*stats)

	return nil
}

func (nr *NodeRegistry) CountDownToClose() {
	log.Println("===> Will close automatically in 60 seconds <===")
	time.Sleep(60 * time.Second)
	nr.mutex.Lock()
	defer nr.mutex.Unlock()

	nr.statusLogger.LogFinished()
	os.Exit(0)
}
