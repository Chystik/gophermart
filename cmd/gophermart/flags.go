package main

import (
	"flag"

	"github.com/Chystik/gophermart/config"
)

func parseFlags(cfg *config.App) {
	// Check interface implementation
	_ = flag.Value(&cfg.Address)
	_ = flag.Value(&cfg.DBuri)
	_ = flag.Value(&cfg.AccrualAddress)

	flag.Var(&cfg.Address, "a", "app address")
	flag.Var(&cfg.DBuri, "d", "database uri")
	flag.Var(&cfg.AccrualAddress, "r", "accal service address")
	flag.Parse()
}
