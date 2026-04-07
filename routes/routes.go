package routes

import (
	"github.com/FoodByMegha/foodbymegha-backend/handlers"
	"github.com/FoodByMegha/foodbymegha-backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// API version 1
	api := r.Group("/api/v1")

	// Public routes — koi bhi access kar sakta hai
	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Protected routes — sirf logged in customer
	protected := api.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/profile", handlers.GetProfile)
		protected.GET("/plans", handlers.GetPlans)
		protected.POST("/subscribe", handlers.Subscribe)
		protected.GET("/my-plan", handlers.GetMyPlan)

		protected.POST("/orders", handlers.CreateOrder)
		protected.GET("/orders", handlers.GetOrders)
		protected.GET("/track/:id", handlers.TrackOrder)
		protected.PUT("/orders/:id/status", handlers.UpdateOrderStatus)
	}
}
