package main

import (
	"log"
	"os"

	"github.com/FoodByMegha/foodbymegha-backend/config"
	"github.com/FoodByMegha/foodbymegha-backend/models"
	"github.com/FoodByMegha/foodbymegha-backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// .env file load karo
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database connect karo
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{}, &models.Plan{}, &models.Subscription{}, &models.Order{}, &models.Payment{})

	// Gin router shuru karo
	r := gin.Default()

	// ✨ CORS add karo
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Sab routes setup karo
	routes.SetupRoutes(r)

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "FoodByMegha backend is live! 🍱",
		})
	})

	// Port lo .env se
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Server start karo
	r.Run(":" + port)
}
