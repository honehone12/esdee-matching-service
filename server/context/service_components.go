package context

import (
	"esdee-matching-service/queue"
	"esdee-matching-service/roller"
)

type ServiceComponents struct {
	queue  queue.QueueHandle
	roller roller.RollerHandle
}

func NewServiceComponents(
	queue queue.QueueHandle,
	roller roller.RollerHandle,
) *ServiceComponents {
	return &ServiceComponents{
		queue:  queue,
		roller: roller,
	}
}

func (c *ServiceComponents) Queue() queue.QueueHandle {
	return c.queue
}

func (c *ServiceComponents) Roller() roller.RollerHandle {
	return c.roller
}
