package status

import (
	"errors"
	"esdee-matching-service/logger"
	"fmt"
	"sync"
	"time"
)

type StatusMapAddHandle interface {
	Count() int64
	Add(string, *Status)
}

type StatusMapRWHandle interface {
	Count() int64
	Item(string) (*Status, error)
}

type StatusMapConsumeHandle interface {
	Count() int64
	Remove(string)
	ReadonlyItem(string) (StatusValue, error)
}

type StatusMapSetting struct {
	MonitoringInterval time.Duration
}

type StatusMap struct {
	setting StatusMapSetting
	mapping sync.Map
	count   int64
	logger  logger.Logger
	eChan   chan error
}

var (
	ErrorNoSuchItem           = errors.New("no such item mapped")
	ErrorUnexpectedTypeMapped = errors.New("mapped interface was not type expected")
)

func NewStatusMap(
	logger logger.Logger,
	monitoringInterval time.Duration,
) *StatusMap {
	return &StatusMap{
		setting: StatusMapSetting{
			MonitoringInterval: monitoringInterval,
		},
		mapping: sync.Map{},
		count:   0,
		logger:  logger,
		eChan:   make(chan error),
	}
}

func (m *StatusMap) Add(uuid string, status *Status) {
	if _, exists := m.mapping.LoadOrStore(uuid, status); !exists {
		m.count++
	}
}

func (m *StatusMap) Remove(uuid string) {
	if _, exists := m.mapping.LoadAndDelete(uuid); exists {
		m.count--
	}
}

func (m *StatusMap) Item(uuid string) (*Status, error) {
	i, ok := m.mapping.Load(uuid)
	if ok {
		status, ok := i.(*Status)
		if !ok {
			return nil, ErrorUnexpectedTypeMapped
		}
		return status, nil
	} else {
		return nil, ErrorNoSuchItem
	}
}

func (m *StatusMap) ReadonlyItem(uuid string) (StatusValue, error) {
	s, err := m.Item(uuid)
	if err != nil {
		return StatusValue{}, err
	}
	return s.value, nil
}

func (m *StatusMap) Count() int64 {
	return m.count
}

func (m *StatusMap) StartMonitoring() <-chan error {
	go m.monitor()
	return m.eChan
}

func (m *StatusMap) monitor() {
	tick := time.Tick(m.setting.MonitoringInterval)
	format := "[MAP] count: %d"
	for range tick {
		if m.count < 0 {
			m.eChan <- fmt.Errorf(format, m.count)
		} else {
			m.logger.Infof(format, m.count)
		}
	}
}
