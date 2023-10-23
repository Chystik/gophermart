package main

import (
	"os"

	"github.com/Chystik/gophermart/config"
	"github.com/joho/godotenv"

	"github.com/caarlos0/env"
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
