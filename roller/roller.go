package roller

import (
	"esdee-matching-service/logger"
	"esdee-matching-service/queue"
	"time"
)

type RollerHandle interface {
}

type RollerSetting struct {
	Interval     time.Duration
	MatchingUnit int64
}

type Roller struct {
	setting RollerSetting
	queue   queue.QueueHandle
	logger  logger.Logger
	eChan   chan error
}

func New(
	logger logger.Logger,
	queue queue.QueueHandle,
	unit int64,
	interval time.Duration,
) *Roller {
	return &Roller{
		setting: RollerSetting{
			Interval:     interval,
			MatchingUnit: unit,
		},
		queue:  queue,
		logger: logger,
		eChan:  make(chan error),
	}
}

func (r *Roller) StartRolling() <-chan error {
	go r.roll()
	return r.eChan
}

func (r *Roller) roll() {
	tick := time.Tick(r.setting.Interval)
	for range tick {
		if r.setting.MatchingUnit <= 0 {
			continue
		}

		len := r.queue.Len()
		numMatched := len / r.setting.MatchingUnit
		r.logger.Infof("%d matched !!", numMatched)
		for i := int64(0); i < numMatched; i++ {

		}
	}
}
