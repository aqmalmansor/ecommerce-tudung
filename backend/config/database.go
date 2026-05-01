package config

import (
	"log"

	"be/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established")
	return db, nil
}

func MigrateDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	log.Println("Database migration completed")
	return nil
}
