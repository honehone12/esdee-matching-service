package handlers

import (
	"esdee-matching-service/server/context"
	"esdee-matching-service/status"
	"esdee-matching-service/ticket"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TicketCreateResponse struct {
	Uuid string
}

func TicketCreate(c echo.Context) error {
	ticket := ticket.New()
	ctx := c.(*context.Context)

	if err := ctx.MatchingEnqueue().Enqueue(ticket); err != nil {
		return err
	}
	ctx.StatusAdd().Add(ticket.String(), status.NewStatus())

	return c.JSON(http.StatusOK, TicketCreateResponse{
		Uuid: ticket.String(),
	})
}
