package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Database DBConfig
}

type DBConfig struct {
	DBAddr string
	DBPort string
	DBName string
	DBUser string
	DBPass string
	DBType string
}

var Config Configuration

func SetupConfig() {
	var raw []byte
	var err error

	if raw, err = os.ReadFile("config.json"); err != nil {
		log.Fatal("Unable to read configuration file: ", err)
	}

	if err = json.Unmarshal(raw, &Config); err != nil {
		log.Fatal("Unable to parse configuration file: ", err)
	}
}
