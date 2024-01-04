package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() (*gorm.DB, error) {
	host := "localhost"
	port := "5432"
	dbName := "orders"
	dbUser := "nawaf"
	password := "n2345%#%3dk"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// DB.AutoMigrate(models.OrderedItem{}, models.Order{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Database connection successful...")

	return DB, nil
}
