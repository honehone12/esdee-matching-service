package server

import (
	"esdee-matching-service/server/context"
	"esdee-matching-service/server/handlers"
	"esdee-matching-service/server/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Server struct {
	echo       *echo.Echo
	metadata   *context.Metadata
	components *context.ServiceComponents
	listenAt   string
}

func NewServer(
	echo *echo.Echo,
	metadata *context.Metadata,
	components *context.ServiceComponents,
	listenAt string,
) *Server {
	return &Server{
		echo:       echo,
		metadata:   metadata,
		components: components,
		listenAt:   listenAt,
	}
}

func (s *Server) ConvertContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := &context.Context{
			Context:           c,
			Metadata:          s.metadata,
			ServiceComponents: s.components,
		}
		return next(ctx)
	}
}

func (s *Server) Run() {
	s.echo.Use(s.ConvertContext)
	s.echo.Validator = validator.NewValidator()
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.Logger())

	s.echo.GET("/", handlers.Root)
	s.echo.GET("/ticket/create", handlers.TicketCreate)
	s.echo.POST("/status/poll", handlers.StatusPoll)
	s.echo.POST("/status/standby", handlers.StatusStandby)

	s.echo.Logger.SetLevel(log.INFO)
	s.echo.Logger.Fatal(s.echo.Start(s.listenAt))
}
