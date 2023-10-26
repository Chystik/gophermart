package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Chystik/gophermart/config"
	"github.com/Chystik/gophermart/run"
)

func main() {
	cfg := config.NewAppConfig()

	err := parseEnv(cfg)
	if err != nil {
		panic(err)
	}
	//parseFlags(cfg)

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	run.App(cfg, quit)
}
