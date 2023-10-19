package handlers

import (
	"errors"
	"esdee-matching-service/server/context"
	"esdee-matching-service/server/errorres"
	"esdee-matching-service/server/form"
	"esdee-matching-service/status"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PollingForm struct {
	Uuid string `form:"uuid" validate:"required,uuid4,min=36,max=36"`
}

func StatusPoll(c echo.Context) error {
	formData := &PollingForm{}
	if err := form.ProcessFormData(c, formData); err != nil {
		return errorres.BadRequest(err, c.Logger())
	}

	mapping := c.(*context.Context).ServiceComponents.StatusConsume()
	current, err := mapping.ReadonlyItem(formData.Uuid)
	if err != nil {
		return errorres.BadRequest(err, c.Logger())
	}

	return c.JSON(http.StatusOK, current)
}

func StatusStandby(c echo.Context) error {
	formData := &PollingForm{}
	if err := form.ProcessFormData(c, formData); err != nil {
		return errorres.BadRequest(err, c.Logger())
	}

	mapping := c.(*context.Context).ServiceComponents.StatusConsume()
	current, err := mapping.ReadonlyItem(formData.Uuid)
	if err != nil {
		return errorres.BadRequest(err, c.Logger())
	}
	if current.StatusCode != status.StatusDone {
		return errorres.BadRequest(
			errors.New("status is still waiting"),
			c.Logger(),
		)
	}

	mapping.Remove(formData.Uuid)
	return c.NoContent(http.StatusOK)
}
