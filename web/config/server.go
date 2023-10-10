package config

import (
	"github.com/Kaikawa1028/go-template/app/errors"
	"os"
	"strconv"
)

type ServerConfig struct {
	Port int
}

func NewServerConfig() (*ServerConfig, error) {
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return nil, errors.Wrap(err, "Could not read the env var 'SERVER_PORT'")
	}

	return &ServerConfig{
		Port: port,
	}, nil
}
