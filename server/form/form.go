package form

import (
	"github.com/labstack/echo/v4"
)

type FormData interface {
}

func ProcessFormData[T FormData](c echo.Context, ptr *T) error {
	if err := c.Bind(ptr); err != nil {
		return err
	}
	if err := c.Validate(ptr); err != nil {
		return err
	}
	return nil
}
