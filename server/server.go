package server

import (
	"esdee-matching-service/server/context"
	"esdee-matching-service/server/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Server struct {
	metadata context.Metadata
	listenAt string
}

func NewServer(name string, version string, listenAt string) *Server {
	return &Server{
		metadata: context.NewBasicMetadata(name, version),
		listenAt: listenAt,
	}
}

func (s *Server) ConvertContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := &context.Context{
			Context:  c,
			Metadata: s.metadata,
		}
		return next(ctx)
	}
}

func (s *Server) Run() {
	e := echo.New()
	e.Use(s.ConvertContext)
	e.Validator = validator.NewValidator()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.Logger.SetLevel(log.WARN)
	e.Logger.Fatal(e.Start(s.listenAt))
}
