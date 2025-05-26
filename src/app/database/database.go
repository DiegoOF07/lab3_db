package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	dsn := "host=database user=postgres password=secret dbname=lab3_products port=5432 sslmode=disable TimeZone=UTC"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}
	return DB
}

func SeedData(db *gorm.DB, path string) error {
	sqlBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error leyendo archivo SQL (%s): %w", path, err)
	}

	sql := string(sqlBytes)

	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("error ejecutando data.sql: %w", err)
	}

	return nil
}
