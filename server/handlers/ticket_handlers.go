package handlers

import (
	"esdee-matching-service/server/context"
	"esdee-matching-service/ticket"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TicketCreateResponse struct {
	Uuid string
}

func TicketCreate(c echo.Context) error {
	ticket := ticket.New()
	if err := c.(*context.Context).Queue().Enqueue(ticket); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, TicketCreateResponse{
		Uuid: ticket.String(),
	})
}
