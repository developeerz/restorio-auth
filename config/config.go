package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ServiceName string
	Access      string // `validate:"required"`
	Refresh     string
	Postgres    string
}

var ConfigService Config

func LoadConfig() {
	err := godotenv.Load()
	log.Info().Msgf("load dotenv: %v", err)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	ConfigService = Config{
		ServiceName: os.Getenv("SERVICE_NAME"),
		Access:      os.Getenv("SECRET_ACCESS"),
		Refresh:     os.Getenv("SECRET_REFRESH"),
		Postgres:    dsn,
	}
}
