package config

import (
	"errors"
	"strconv"
	"strings"
)

type (
	App struct {
		Address        `env:"RUN_ADDRESS"`
		DBuri          `env:"DATABASE_URI"`
		AccrualAddress `env:"ACCRUAL_SYSTEM_ADDRESS"`
		JWTkey         []byte
	}

	Address        string
	DBuri          string
	AccrualAddress string
)

func NewAppConfig() *App {
	return &App{
		Address: ":8080",
		JWTkey:  []byte("my_secret_key"),
	}
}

func (addr Address) String() string {
	return string(addr)
}

func (addr *Address) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("expect address in a form host:port")
	}
	_, err := strconv.Atoi(hp[1])
	if err != nil {
		return errors.New("only digits allowed for port in a form host:port")
	}
	*addr = Address(s)
	return nil
}

func (db DBuri) String() string {
	return string(db)
}

func (db *DBuri) Set(s string) error {
	*db = DBuri(s)
	return nil
}

func (accr AccrualAddress) String() string {
	return string(accr)
}

func (accr *AccrualAddress) Set(s string) error {
	*accr = AccrualAddress(s)
	return nil
}
