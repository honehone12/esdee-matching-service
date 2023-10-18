package matching

import (
	"errors"
	"esdee-matching-service/logger"
	"esdee-matching-service/status"
	"time"
)

type RollerHandle interface {
	// use this handle for example change interval etc...
}

type RollerSetting struct {
	Interval     time.Duration
	MatchingUnit int64
}

type Roller struct {
	setting RollerSetting
	queue   MatchingDequeueHandle
	mapping status.StatusMapRWHandle
	logger  logger.Logger
	eChan   chan error
}

var (
	ErrorInvalidMatchingUnit = errors.New("matching unit is set as 0 or smaller")
)

func NewRoller(
	logger logger.Logger,
	queue MatchingDequeueHandle,
	mapping status.StatusMapRWHandle,
	unit int64,
	interval time.Duration,
) *Roller {
	return &Roller{
		setting: RollerSetting{
			Interval:     interval,
			MatchingUnit: unit,
		},
		queue:   queue,
		mapping: mapping,
		logger:  logger,
		eChan:   make(chan error),
	}
}

func (r *Roller) StartRolling() <-chan error {
	go r.roll()
	return r.eChan
}

func (r *Roller) roll() {
	tick := time.Tick(r.setting.Interval)
TICK:
	for range tick {
		if r.setting.MatchingUnit <= 0 {
			r.eChan <- ErrorInvalidMatchingUnit
			continue
		}

		len := r.queue.Len()
		numMatched := len / r.setting.MatchingUnit
		if numMatched == 0 {
			continue
		}

		for i := int64(0); i < numMatched; i++ {
			tickets, err := r.queue.Dequeue(r.setting.MatchingUnit)
			if err != nil {
				r.eChan <- err
				continue TICK
			}

			ip := "127.0.0.1:9999"
			uuids := make([]string, r.setting.MatchingUnit)
			for j := int64(0); j < r.setting.MatchingUnit; j++ {
				uuids[j] = tickets[j].String()
			}

			for j := int64(0); j < r.setting.MatchingUnit; j++ {
				s, err := r.mapping.Item(tickets[j].String())
				if err != nil {
					r.eChan <- err
					continue TICK
				}

				if err = s.SetAsDone(ip, uuids); err != nil {
					r.eChan <- err
					continue TICK
				}
			}
			r.logger.Infof("[MATCHING] ip: %s, uuids: %v", ip, uuids)
		}
	}
}
