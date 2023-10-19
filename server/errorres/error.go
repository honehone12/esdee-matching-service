package errorres

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func BadRequest(err error, logger echo.Logger) error {
	logger.Warn(err)
	return echo.NewHTTPError(
		http.StatusBadRequest,
		"input value is not valid",
	)
}

func ServiceError(err error, logger echo.Logger) error {
	logger.Error(err)
	return echo.NewHTTPError(
		http.StatusInternalServerError,
		"service has unexpected error",
	)
}
