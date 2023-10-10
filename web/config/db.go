package config

import (
	"github.com/Kaikawa1028/go-template/app/errors"
	"os"
	"strconv"
)

type DBConfig struct {
	User     string
	Password string
	Database string
	Host     string
	Port     int
	TimeZone string
}

func NewDBConfig() (*DBConfig, error) {
	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		return nil, errors.Wrap(err, "Could not read the env var 'MYSQL_PORT'")
	}

	dbConfig := &DBConfig{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DATABASE"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     port,
		TimeZone: os.Getenv("TZ"),
	}

	return dbConfig, nil
}
