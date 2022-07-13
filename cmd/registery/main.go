package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/rpc"

	"github.com/korkmazkadir/ida-gossip/registery"
)

const configFile = "config.json"

func main() {

	statusLogger := registery.NewStatusLogger()
	defer func() {
		if r := recover(); r != nil {
			statusLogger.LogFailed()
		}
	}()

	nodeConfig := readConfigFromFile()

	nodeRegistry := registery.NewNodeRegistry(nodeConfig, statusLogger)

	err := rpc.Register(nodeRegistry)
	if err != nil {
		panic(err)
	}

	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	log.Printf("registery service started and listening on :1234\n")

	for {
		conn, _ := l.Accept()
		go func() {
			rpc.ServeConn(conn)
		}()
	}
}

func readConfigFromFile() registery.NodeConfig {

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	config := registery.NodeConfig{}
	json.Unmarshal(data, &config)

	return config
}
