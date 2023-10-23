package main

import (
	"flag"

	"github.com/Chystik/gophermart/config"
)

func parseFlags(cfg *config.App) {
	flag.StringVar(&cfg.Address, "a", "", "app address")
	flag.StringVar(&cfg.DBuri, "d", "", "database uri")
	flag.StringVar(&cfg.AccrualAddress, "r", "", "accal service address")
	flag.Parse()
}
