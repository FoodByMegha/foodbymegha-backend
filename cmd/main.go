package main

import (
	"os"

	"github.com/FoodByMegha/foodbymegha-backend/config"
	"github.com/FoodByMegha/foodbymegha-backend/cron"
	"github.com/FoodByMegha/foodbymegha-backend/models"
	"github.com/FoodByMegha/foodbymegha-backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{}, &models.Plan{}, &models.Subscription{}, &models.Order{}, &models.Payment{})

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://foodbymegha-frontend.vercel.app",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"}, // ✅ PATCH add hua
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	routes.SetupRoutes(r)

	cron.StartCronJobs() // ✅ Cron job shuru

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "FoodByMegha backend is live! 🍱",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
