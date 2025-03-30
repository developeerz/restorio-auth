package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Access   string
	Refresh  string
	Postgres string
}

var ConfigService Config

func LoadConfig() {
	godotenv.Load()

	ConfigService = Config{
		Access:   os.Getenv("SECRET_ACCESS"),
		Refresh:  os.Getenv("SECRET_REFRESH"),
		Postgres: os.Getenv("POSTGRES_CONFIG"),
	}
}
