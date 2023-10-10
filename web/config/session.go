package config

import (
	"os"
	"strconv"

	"github.com/Kaikawa1028/go-template/app/errors"
)

type SessionConfig struct {
	Time int
}

func NewSessionConfig() (*SessionConfig, error) {
	time, err := strconv.Atoi(os.Getenv("SESSION_TIME"))
	if err != nil {
		return nil, errors.Wrap(err, "Could not read the env var 'SESSION_TIME'")
	}

	return &SessionConfig{
		Time: time,
	}, nil
}
