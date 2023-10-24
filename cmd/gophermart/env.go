package main

import (
	"os"

	"github.com/Chystik/gophermart/config"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

func parseEnv(cfg *config.App, dotEnvFile string) error {
	if _, err := os.Stat(dotEnvFile); err == nil {
		errLoad := godotenv.Load(dotEnvFile)
		if errLoad != nil {
			return errLoad
		}
	}
	return env.Parse(cfg)
}
