package common

import (
	"fmt"
	"log"
)

type EventType int

const (
	FirstChunkReceived EventType = iota
	MessageReceived
	NetworkUsage
	DisseminationFailure
	MeanSendTime
)

func (e EventType) String() string {
	switch e {
	case FirstChunkReceived:
		return "FIRST_CHUNK_RECEIVED"
	case MessageReceived:
		return "MESSAGE_RECEIVED"
	case NetworkUsage:
		return "NETWORK_USAGE"
	case DisseminationFailure:
		return "DISSEMINATION_FAILURE"
	case MeanSendTime:
		return "MEAN_SEND_TIME"
	default:
		panic(fmt.Errorf("undefined enum value %d", e))
	}
}

type Event struct {
	Round        int
	Type         EventType
	ElapsedTime  int
	NetworkUsage int64
	MeanSendTime float64
}

type StatList struct {
	IPAddress  string
	PortNumber int
	NodeID     int
	Events     []Event
}

type StatLogger struct {
	nodeID int
	events []Event
}

func NewStatLogger(nodeID int) *StatLogger {
	return &StatLogger{nodeID: nodeID}
}

func (s *StatLogger) FirstChunkReceived(round int, elapsedTime int64) {
	log.Printf("stats\t%d\t%d\t%s\t%d\t", s.nodeID, round, "FIRST_CHUNK_RECEIVED", elapsedTime)
	s.events = append(s.events, Event{Round: round, Type: FirstChunkReceived, ElapsedTime: int(elapsedTime)})
}

func (s *StatLogger) MessageReceived(round int, elapsedTime int64) {
	log.Printf("stats\t%d\t%d\t%s\t%d\t", s.nodeID, round, "MESSAGE_RECEIVED", elapsedTime)
	s.events = append(s.events, Event{Round: round, Type: MessageReceived, ElapsedTime: int(elapsedTime)})
}

func (s *StatLogger) NetworkUsage(round int, usage int64) {
	log.Printf("stats\t%d\t%d\t%s\t%d\t", s.nodeID, round, "NETWORK_USAGE", usage)
	s.events = append(s.events, Event{Round: round, Type: NetworkUsage, NetworkUsage: usage})
}

func (s *StatLogger) DisseminationFailure(round int, leader int) {
	log.Printf("stats\t%d\t%d\t%s\t%d\t", s.nodeID, round, "DISSEMINATION_FAILURE", leader)
	s.events = append(s.events, Event{Round: round, Type: DisseminationFailure, ElapsedTime: leader})
}

func (s *StatLogger) MeanSendTime(round int, time float64) {
	log.Printf("stats\t%d\t%d\t%s\t%f\t", s.nodeID, round, "MEAN_SEND_TIME", time)
	s.events = append(s.events, Event{Round: round, Type: MeanSendTime, MeanSendTime: time})
}

func (s *StatLogger) GetEvents() []Event {
	return s.events
}
