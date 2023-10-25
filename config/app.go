package config

type App struct {
	Address        string `env:"RUN_ADDRESS"`
	DBuri          string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTkey         []byte
}

func NewAppConfig() *App {
	return &App{
		JWTkey: []byte("my_secret_key"),
	}
}
