package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/korkmazkadir/ida-gossip/common"
	"github.com/korkmazkadir/ida-gossip/dissemination"
	"github.com/korkmazkadir/ida-gossip/network"
	"github.com/korkmazkadir/ida-gossip/registery"
)

func main() {

	hostname := getEnvWithDefault("NODE_HOSTNAME", "127.0.0.1")
	registryAddress := getEnvWithDefault("REGISTRY_ADDRESS", "localhost:1234")
	processIndex := getEnvWithDefault("PROCESS_INDEX", "-1")

	log.Printf("Process Index: %s\n", processIndex)

	demux := common.NewDemultiplexer(0)
	server := network.NewServer(demux)

	err := rpc.Register(server)
	if err != nil {
		panic(err)
	}

	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf("%s:", hostname))
	if e != nil {
		log.Fatal("listen error:", e)
	}

	// start serving
	go func() {
		for {
			conn, _ := l.Accept()
			go func() {
				rpc.ServeConn(conn)
			}()
		}
	}()

	log.Printf("p2p server started on %s\n", l.Addr().String())
	nodeInfo := getNodeInfo(l.Addr().String())

	registry := registery.NewRegistryClient(registryAddress, nodeInfo)

	nodeInfo.ID = registry.RegisterNode()
	log.Printf("node registeration successful, assigned ID is %d\n", nodeInfo.ID)

	nodeConfig := registry.GetConfig()

	var nodeList []registery.NodeInfo

	for {
		nodeList = registry.GetNodeList()
		nodeCount := len(nodeList)
		if nodeCount == nodeConfig.NodeCount {
			break
		}
		time.Sleep(2 * time.Second)
		log.Printf("received node list %d/%d\n", nodeCount, nodeConfig.NodeCount)
	}

	peerSet := createPeerSet(nodeList, nodeConfig.GossipFanout, nodeInfo.ID, nodeInfo.IPAddress, nodeConfig.ConnectionCount)
	statLogger := common.NewStatLogger(nodeInfo.ID)
	rapidchain := dissemination.NewDisseminator(demux, nodeConfig, peerSet, statLogger)

	runConsensus(rapidchain, nodeConfig.EndRound, nodeConfig.RoundSleepTime, nodeInfo.ID, nodeConfig.SourceCount, nodeConfig.MessageSize, nodeList)

	// collects stats abd uploads to registry
	log.Printf("uploading stats to the registry\n")
	events := statLogger.GetEvents()
	statList := common.StatList{IPAddress: nodeInfo.IPAddress, PortNumber: nodeInfo.PortNumber, NodeID: nodeInfo.ID, Events: events}
	registry.UploadStats(statList)

	log.Printf("reached target round count. Shutting down in 1 minute\n")
	time.Sleep(1 * time.Minute)

	log.Printf("exiting as expected...\n")
}

func createPeerSet(nodeList []registery.NodeInfo, fanOut int, nodeID int, localIPAddress string, connectionCount int) network.PeerSet {

	var copyNodeList []registery.NodeInfo
	copyNodeList = append(copyNodeList, nodeList...)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(copyNodeList), func(i, j int) { copyNodeList[i], copyNodeList[j] = copyNodeList[j], copyNodeList[i] })

	peerSet := network.PeerSet{}

	peerCount := 0
	for i := 0; i < len(copyNodeList); i++ {

		if peerCount == fanOut {
			break
		}

		peer := copyNodeList[i]
		//TODO: do not connect nodes from local machine
		if peer.ID == nodeID || peer.IPAddress == localIPAddress {
			continue
		}

		err := peerSet.AddPeer(peer.IPAddress, peer.PortNumber, connectionCount)
		if err != nil {
			panic(err)
		}
		log.Printf("new peer added: %s:%d ID %d\n", peer.IPAddress, peer.PortNumber, peer.ID)
		peerCount++
	}

	return peerSet
}

func getNodeInfo(netAddress string) registery.NodeInfo {
	tokens := strings.Split(netAddress, ":")

	ipAddress := tokens[0]
	portNumber, err := strconv.Atoi(tokens[1])
	if err != nil {
		panic(err)
	}

	return registery.NodeInfo{IPAddress: ipAddress, PortNumber: portNumber}
}

func runConsensus(rc *dissemination.Disseminator, numberOfRounds int, roundSleepTime int, nodeID int, leaderCount int, blockSize int, nodeList []registery.NodeInfo) {

	currentRound := 1
	for currentRound <= numberOfRounds {

		log.Printf("+++++++++ Round %d +++++++++++++++\n", currentRound)

		var messages []common.Message

		// if elected as a leader submits a block
		isElected, electedLeaders := isElectedAsLeader(nodeList, currentRound, nodeID, leaderCount)

		if isElected {
			log.Println("elected as leader")
			b := createBlock(currentRound, nodeID, blockSize, leaderCount)
			rc.SubmitMessage(currentRound, b)
		}

		// TODO: is it better to log individual messages?
		// waits to deliver the block
		log.Printf("waiting to deliver messages...\n")
		messages, unresponsiveLeaders := rc.WaitForMessage(currentRound, electedLeaders)

		if unresponsiveLeaders != nil {
			log.Printf("Unresponsive leader detected: %v\n", unresponsiveLeaders)
			nodeList = removeUnresponsiveLeaders(unresponsiveLeaders, nodeList)
			log.Printf("Unresponsive leaders are removed.")
			continue
		}

		log.Printf("all messages delivered.\n")
		payloadSize := 0
		for i := range messages {
			log.Printf("Round: %d Message[%d] %x\n", currentRound, i, common.EncodeBase64(messages[i].Hash())[:15])
			payloadSize += len(messages[i].Payload)
		}

		log.Printf("round finished, payload size payload size: %d bytes\n", payloadSize)

		currentRound++

		sleepTime := int64(roundSleepTime*1000) - (time.Now().UnixMilli() - messages[0].Time)
		log.Printf("sleeping for %d ms\n", sleepTime)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)

		//sleepTime := time.Duration(roundSleepTime) * time.Second
		//log.Printf("sleeping for %s\n", sleepTime)
		//time.Sleep(sleepTime)

	}

}

////////////////////
///// utils ////////
////////////////////

func removeUnresponsiveLeaders(unresponsiveLeaders []int, nodeList []registery.NodeInfo) []registery.NodeInfo {

	var newNodeList []registery.NodeInfo
	for _, nodeInfo := range nodeList {

		isInList := false
		for _, l := range unresponsiveLeaders {
			if l == nodeInfo.ID {
				isInList = true
				break
			}
		}

		if isInList == false {
			newNodeList = append(newNodeList, nodeInfo)
		}

	}

	return newNodeList
}

func createBlock(round int, nodeID int, blockSize int, leaderCount int) common.Message {

	payloadSize := blockSize

	block := common.Message{
		Round:   round,
		Issuer:  nodeID,
		Time:    time.Now().UnixMilli(),
		Payload: common.GetRandomByteSlice(payloadSize),
	}

	return block
}

func getEnvWithDefault(key string, defaultValue string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		val = defaultValue
	}

	log.Printf("%s=%s\n", key, val)
	return val
}

func isElectedAsLeader(nodeList []registery.NodeInfo, round int, nodeID int, leaderCount int) (bool, []int) {

	// assumes that node list is same for all nodes
	// shuffle the node list using round number as the source of randomness
	rand.Seed(int64(round))
	rand.Shuffle(len(nodeList), func(i, j int) { nodeList[i], nodeList[j] = nodeList[j], nodeList[i] })

	var electedLeaders []int
	isElected := false
	for i := 0; i < leaderCount; i++ {
		electedLeaders = append(electedLeaders, nodeList[i].ID)
		if nodeList[i].ID == nodeID {
			log.Println("=== elects as a leader ===")
			isElected = true
		}
	}

	log.Printf("Elected leaders: %v\n", electedLeaders)

	return isElected, electedLeaders
}
