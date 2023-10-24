package matching

import (
	"errors"
	"esdee-matching-service/game"
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
	gameIp  game.GameServerIp
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
	gameIp game.GameServerIp,
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
		gameIp:  gameIp,
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

			address, port := r.gameIp.GetNextProcessIp()
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

				if err = s.SetAsDone(address, port, uuids); err != nil {
					r.eChan <- err
					continue TICK
				}
			}
			r.logger.Infof(
				"[MATCHING] address: %s, port: %s, uuids: %v",
				address, port, uuids,
			)
		}
	}
}
