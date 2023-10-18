package matching

import (
	"errors"
	"esdee-matching-service/logger"
	"esdee-matching-service/ticket"
	"time"

	lib "github.com/Workiva/go-datastructures/queue"
)

type MatchingEnqueueHandle interface {
	Len() int64
	Enqueue(*ticket.Ticket) error
}

type MatchingDequeueHandle interface {
	Len() int64
	Dequeue(int64) ([]*ticket.Ticket, error)
}

type MatchingQueueSetting struct {
	InitialCapacity    int64
	MonitoringInterval time.Duration
}

type MatchingQueue struct {
	setting MatchingQueueSetting
	queue   *lib.Queue
	logger  logger.Logger
	eChan   chan error
}

var (
	ErrorUnexpectedTypeQueued = errors.New("queued interface was not type expected")
)

func NewMatchingQueue(
	logger logger.Logger,
	initialCapacity int64,
	monitoringInterval time.Duration,
) *MatchingQueue {
	return &MatchingQueue{
		setting: MatchingQueueSetting{
			InitialCapacity:    initialCapacity,
			MonitoringInterval: monitoringInterval,
		},
		queue:  lib.New(initialCapacity),
		logger: logger,
		eChan:  make(chan error),
	}
}

func (q *MatchingQueue) Len() int64 {
	return q.queue.Len()
}

func (q *MatchingQueue) Enqueue(ticket *ticket.Ticket) error {
	return q.queue.Put(ticket)
}

func (q *MatchingQueue) Dequeue(n int64) ([]*ticket.Ticket, error) {
	interfaces, err := q.queue.Get(n)
	if err != nil {
		return nil, err
	}

	len := len(interfaces)
	tickets := make([]*ticket.Ticket, len)
	for i := 0; i < len; i++ {
		t, ok := interfaces[i].(*ticket.Ticket)
		if !ok {
			return nil, ErrorUnexpectedTypeQueued
		}
		tickets[i] = t
	}
	return tickets, nil
}

func (q *MatchingQueue) StartMonitoring() <-chan error {
	go q.monitor()
	return q.eChan
}

func (q *MatchingQueue) monitor() {
	tick := time.Tick(q.setting.MonitoringInterval)
	format := "[QUEUE] length: %d"
	for range tick {
		len := q.queue.Len()
		q.logger.Infof(format, len)
	}
}
