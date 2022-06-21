package common

import (
	"fmt"
	"log"
)

type EventType int

const (
	FirstChunkReceived EventType = iota
	MessageReceived
	QueueLength
)

func (e EventType) String() string {
	switch e {
	case FirstChunkReceived:
		return "FIRST_CHUNK_RECEIVED"
	case MessageReceived:
		return "MESSAGE_RECEIVED"
	case QueueLength:
		return "QUEUE_LENGTH"
	default:
		panic(fmt.Errorf("undefined enum value %d", e))
	}
}

type Event struct {
	Round       int
	Type        EventType
	ElapsedTime int
	QueuLength  float64
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

func (s *StatLogger) AvgQueuLength(round int, queueLength float64) {
	log.Printf("stats\t%d\t%d\t%s\t%f\t", s.nodeID, round, "AVG_QUEUE_LENGTH", queueLength)
	s.events = append(s.events, Event{Round: round, Type: QueueLength, QueuLength: queueLength})
}

func (s *StatLogger) GetEvents() []Event {
	return s.events
}
