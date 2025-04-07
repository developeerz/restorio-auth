package config

import (
	"fmt"
	"log"
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
	err := godotenv.Load()
	log.Println("LoadConfig: ", err)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	ConfigService = Config{
		Access:   os.Getenv("SECRET_ACCESS"),
		Refresh:  os.Getenv("SECRET_REFRESH"),
		Postgres: dsn,
	}
}
