package database

import (
	"fmt"
	"log"

	"github.com/developeerz/restorio-auth/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := config.ConfigService.Postgres
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	DB = db
	fmt.Println("Database connected successfully")
}
