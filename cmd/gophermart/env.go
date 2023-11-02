package main

import (
	"os"

	"github.com/Chystik/gophermart/config"
)

func parseEnv(cfg *config.App) {
	cfg.Address = config.Address(os.Getenv("RUN_ADDRESS"))
	cfg.DBuri = config.DBuri(os.Getenv("DATABASE_URI"))
	cfg.AccrualAddress = config.AccrualAddress(os.Getenv("ACCRUAL_SYSTEM_ADDRESS"))
}
