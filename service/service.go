package service

import (
	"esdee-matching-service/matching"
	"esdee-matching-service/server"
	"esdee-matching-service/server/context"
	"esdee-matching-service/status"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func ErrorEmptyEnvParam(env string) error {
	return fmt.Errorf("env param %s is empty", env)
}

func parseParams() (string, string, string, error) {
	if err := godotenv.Load(); err != nil {
		return "", "", "", err
	}

	env := "ESDEE_SERVICE_NAME"
	name := os.Getenv(env)
	if len(name) == 0 {
		return "", "", "", ErrorEmptyEnvParam(env)
	}

	env = "ESDEE_VERSION"
	version := os.Getenv(env)
	if len(version) == 0 {
		return "", "", "", ErrorEmptyEnvParam(env)
	}

	env = "ESDEE_SERVER_LISTEN_AT"
	at := os.Getenv(env)
	if len(at) == 0 {
		return "", "", "", ErrorEmptyEnvParam(env)
	}

	return name, version, at, nil
}

func Run() {
	name, version, listenAt, err := parseParams()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	q := matching.NewMatchingQueue(e.Logger, 100, 10*time.Second)
	m := status.NewStatusMap(e.Logger, 10*time.Second)
	r := matching.NewRoller(e.Logger, q, m, 2, 100*time.Millisecond)

	go Catch(
		e.Logger,
		q.StartMonitoring(),
		m.StartMonitoring(),
		r.StartRolling(),
	)

	server.NewServer(
		e,
		context.NewMetadata(name, version),
		context.NewServiceComponents(q, r, m, m),
		listenAt,
	).Run()
}

func Catch(
	logger echo.Logger,
	qErr <-chan error,
	mErr <-chan error,
	rErr <-chan error,
) {
	var err error
	select {
	case err = <-qErr:
	case err = <-mErr:
	case err = <-rErr:
	}

	// fatal here for Debugging
	logger.Fatal(err)
}
