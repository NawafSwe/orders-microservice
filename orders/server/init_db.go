package main

import (
	"fmt"

	"github.com/nawafswe/orders-service/orders/server/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() (*gorm.DB, error) {
	host := "localhost"
	port := "5432"
	dbName := "orders"
	dbUser := "orders_admin"
	password := "fish0r3$fkds"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(models.Order{}, models.OrderedItem{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Database connection successful...")
	return DB, nil
}
