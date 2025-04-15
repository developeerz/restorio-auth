package database

import (
	"fmt"

	"github.com/developeerz/restorio-auth/config"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgresConnect() (*gorm.DB, error) {
	dsn := config.ConfigService.Postgres

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}

	log.Info().Msg("database connected successfully")

	return db, nil
}
