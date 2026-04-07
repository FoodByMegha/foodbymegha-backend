package config

import (
	"fmt"
	"log"
	"os"

	"github.com/FoodByMegha/foodbymegha-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// .env se database details lo
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Connection string banao
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// Database se connect karo
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database se connect nahi ho paya!", err)
	}

	// Automatically tables banao
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("AutoMigrate fail ho gaya!", err)
	}

	DB = db
	log.Println("Database connected aur tables ready! 🎉")
}
