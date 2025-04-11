package database

import (
	"log"

	"github.com/developeerz/restorio-auth/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := config.ConfigService.Postgres

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}

	log.Println("database connected successfully")

	return db
}
