package service

import (
	"esdee-matching-service/logger"
	"esdee-matching-service/queue"
	"esdee-matching-service/roller"
	"esdee-matching-service/server"
	"esdee-matching-service/server/context"
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
	q := queue.New(e.Logger, 100, 10*time.Second)
	r := roller.New(e.Logger, q, 2, 100*time.Millisecond)
	qErr := q.StartMonitoring()
	rErr := r.StartRolling()
	go Catch(qErr, e.Logger)
	go Catch(rErr, e.Logger)

	server.NewServer(
		e,
		context.NewMetadata(name, version),
		context.NewServiceComponents(q, r),
		listenAt,
	).Run()
}

func Catch(eChan <-chan error, logger logger.Logger) {
	for err := range eChan {
		logger.Error(err)
	}
}
