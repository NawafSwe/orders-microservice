package db

import (
	"fmt"
	"github.com/nawafswe/orders-service/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err := DB.AutoMigrate(models.Order{}, models.OrderedItem{}); err != nil {
		log.Fatal("failed to migrate db tables, err: %w", err)
	}
	if err != nil {
		return nil, err
	}
	fmt.Println("Database connection successful...")
	return DB, nil
}
