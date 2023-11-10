package conf

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"8080"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
	DbUrl    string `env:"DB_URL"`
	GinMode  string `env:"GIN_MODE" envDefault:"release"`

	RedisAddr string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPass string `env:"REDIS_PASS" envDefault:""`
	RedisDB   int    `env:"REDIS_DB" envDefault:"0"`
}

func NewConfig() Config {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := Config{}

	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
