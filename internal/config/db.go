package config

import (
	"log"
	"os"
	"tahap2/internal/domain"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	connDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db, err := connDB.DB()
	if err != nil {
		log.Panicf("failed to get database: %v", err)
	}
	db.SetMaxOpenConns(500)
	db.SetMaxIdleConns(125)
	db.SetConnMaxLifetime(15 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Panicf("failed to ping database: %v", err)
	}

	// auto migrate models
	err = connDB.AutoMigrate(&domain.User{}, &domain.Transaction{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return connDB
}
