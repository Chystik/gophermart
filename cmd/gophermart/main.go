package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Chystik/gophermart/config"
	"github.com/Chystik/gophermart/run"

	"github.com/joho/godotenv"
)

func main() {
	cfg := config.NewAppConfig()

	if osEnv := os.Getenv("ENVIRONMENT"); osEnv == "dev" {
		err := godotenv.Load(".env.dev")
		if err != nil {
			panic(err)
		}
	}

	parseEnv(cfg)
	parseFlags(cfg)

	// channel for Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	run.App(cfg, quit)
}
