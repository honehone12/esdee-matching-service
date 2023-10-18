package status

import (
	"errors"
	"sync"
)

const (
	StatusDone    = 1
	StatusWaiting = 2
)

type MatchingResult struct {
	Ip    string
	Uuids []string
}

type StatusValue struct {
	StatusCode int
	Result     MatchingResult
}

type Status struct {
	value  StatusValue
	rwLock sync.RWMutex
}

var (
	ErrorStatusCodeAlreadyDone = errors.New("status code is already set as done")
)

func NewStatus() *Status {
	return &Status{
		value: StatusValue{
			StatusCode: StatusWaiting,
			Result: MatchingResult{
				Ip:    "",
				Uuids: nil,
			},
		},
		rwLock: sync.RWMutex{},
	}
}

func (s *Status) IsDone() bool {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	return s.value.StatusCode == StatusDone
}

func (s *Status) SetAsDone(ip string, uuids []string) error {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	if s.value.StatusCode == StatusDone {
		return ErrorStatusCodeAlreadyDone
	}

	s.value.Result.Ip = ip
	s.value.Result.Uuids = uuids
	s.value.StatusCode = StatusDone
	return nil
}
