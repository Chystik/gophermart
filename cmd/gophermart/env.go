package main

import (
	"os"

	"github.com/Chystik/gophermart/config"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

func parseEnv(cfg *config.App) error {
	if osEnv := os.Getenv("ENVIRONMENT"); osEnv == "dev" {
		err := godotenv.Load(".env.dev")
		if err != nil {
			return err
		}
	}
	return env.Parse(cfg)
}
