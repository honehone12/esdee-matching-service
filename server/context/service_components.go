package context

import (
	"esdee-matching-service/matching"
	"esdee-matching-service/status"
)

type ServiceComponents struct {
	matchingEnqueue matching.MatchingEnqueueHandle
	roller          matching.RollerHandle
	statusAdd       status.StatusMapAddHandle
	statusConsume   status.StatusMapConsumeHandle
}

func NewServiceComponents(
	enqueue matching.MatchingEnqueueHandle,
	roller matching.RollerHandle,
	add status.StatusMapAddHandle,
	consume status.StatusMapConsumeHandle,
) *ServiceComponents {
	return &ServiceComponents{
		matchingEnqueue: enqueue,
		roller:          roller,
		statusAdd:       add,
		statusConsume:   consume,
	}
}

func (c *ServiceComponents) MatchingEnqueue() matching.MatchingEnqueueHandle {
	return c.matchingEnqueue
}

func (c *ServiceComponents) Roller() matching.RollerHandle {
	return c.roller
}

func (c *ServiceComponents) StatusAdd() status.StatusMapAddHandle {
	return c.statusAdd
}

func (c *ServiceComponents) StatusConsume() status.StatusMapConsumeHandle {
	return c.statusConsume
}
