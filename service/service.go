package service

import (
	"esdee-matching-service/server"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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

	s := server.NewServer(name, version, listenAt)
	s.Run()
}
