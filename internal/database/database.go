package database

import (
	"fmt"
	"log"

	"github.com/adipras/api-gladiatore/internal/config"
	"github.com/adipras/api-gladiatore/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func New(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DatabaseUser,
		cfg.DatabasePass,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// If database doesn't exist, create it
		dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DatabaseUser,
			cfg.DatabasePass,
			cfg.DatabaseHost,
			cfg.DatabasePort,
		)
		
		dbTemp, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
		}
		
		// Create database
		err = dbTemp.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cfg.DatabaseName)).Error
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %v", err)
		}
		
		sqlDB, _ := dbTemp.DB()
		sqlDB.Close()
		
		// Try connecting again with the database
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database after creation: %v", err)
		}
	}

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.Service{},
		&models.Endpoint{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return &Database{db}, nil
}