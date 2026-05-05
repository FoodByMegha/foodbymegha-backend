package routes

import (
	"github.com/FoodByMegha/foodbymegha-backend/handlers"
	"github.com/FoodByMegha/foodbymegha-backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	protected := api.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/profile", handlers.GetProfile)
		protected.GET("/plans", handlers.GetPlans)
		protected.GET("/my-plan", handlers.GetMyPlan)

		protected.POST("/orders", handlers.CreateOrder)
		protected.GET("/orders", handlers.GetOrders)
		protected.GET("/track/:id", handlers.TrackOrder)
		protected.PATCH("/orders/:id/note", handlers.UpdateOrderNote) // ✅ Naya

		protected.POST("/payment", handlers.CreatePayment)
		protected.POST("/payment/verify", handlers.VerifyPayment)
		protected.GET("/payment/history", handlers.GetPaymentHistory)
	}

	admin := api.Group("/admin")
	admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		admin.GET("/orders", handlers.GetAllOrders)
		admin.POST("/menu", handlers.AddMenu)
		admin.GET("/revenue", handlers.GetRevenue)
		admin.PUT("/orders/:id/status", handlers.UpdateOrderStatus)
	}
}
