package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Chystik/gophermart/config"
	"github.com/Chystik/gophermart/run"
)

const dotEnvFile string = ".env.dev"

func main() {
	cfg := config.NewAppConfig()

	parseFlags(cfg)
	fmt.Printf("FLAGS: %#v\n", cfg)
	err := parseEnv(cfg, dotEnvFile)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("ENV: %#v\n", cfg)

	// Graceful shutdown setup
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	run.App(cfg, quit)
}
