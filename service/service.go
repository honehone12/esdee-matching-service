package service

import (
	"esdee-matching-service/game"
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

const (
	DummyGameAddress = "127.0.0.1"
	DummyGamePort    = "9999"
)

const (
	MatchingUnit           = 2
	MatchingRollerInterval = 100 * time.Millisecond

	MatchingQueueInitialCapacity    = 100
	MatchingQueueMonitoringInterval = 10 * time.Second
	StatusMapMonitoringInterval     = 10 * time.Second
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
	g := game.NewDummy(DummyGameAddress, DummyGamePort)
	q := matching.NewMatchingQueue(
		e.Logger,
		MatchingQueueInitialCapacity,
		MatchingQueueMonitoringInterval,
	)
	m := status.NewStatusMap(e.Logger, StatusMapMonitoringInterval)
	r := matching.NewRoller(
		e.Logger, q, m, g,
		MatchingUnit,
		MatchingRollerInterval,
	)

	go Catch(
		e.Logger,
		q.StartMonitoring(),
		m.StartMonitoring(),
		r.StartRolling(),
	)

	server.NewServer(
		e,
		context.NewMetadata(name, version),
		context.NewServiceComponents(q, r, m, m, g),
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
