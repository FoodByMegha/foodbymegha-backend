package handlers

import (
	"net/http"
	"time"

	"github.com/FoodByMegha/foodbymegha-backend/config"
	"github.com/FoodByMegha/foodbymegha-backend/models"
	"github.com/gin-gonic/gin"
)

// GET /admin/orders — saare orders dekho
func GetAllOrders(c *gin.Context) {
	var orders []models.Order
	config.DB.Preload("User").Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"total_orders": len(orders),
		"orders":       orders,
	})
}

// POST /admin/menu — naya plan/menu add karo
func AddMenu(c *gin.Context) {
	var input struct {
		Name         string  `json:"name" binding:"required"`
		Description  string  `json:"description"`
		Price        float64 `json:"price" binding:"required"`
		DurationDays int     `json:"duration_days" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sab fields bharo"})
		return
	}

	plan := models.Plan{
		Name:         input.Name,
		Description:  input.Description,
		Price:        input.Price,
		DurationDays: input.DurationDays,
		IsActive:     true,
	}
	config.DB.Create(&plan)

	c.JSON(http.StatusOK, gin.H{
		"message": "Naya plan add ho gaya! 🍱",
		"plan":    plan,
	})
}

// GET /admin/revenue — aaj ki kamaai
func GetRevenue(c *gin.Context) {
	// Aaj ki date
	today := time.Now().Format("2006-01-02")

	// Aaj ke paid payments
	var payments []models.Payment
	config.DB.Where("status = ? AND DATE(created_at) = ?", "paid", today).Find(&payments)

	// Total calculate karo
	var totalRevenue float64
	for _, p := range payments {
		totalRevenue += p.Amount
	}

	// Total subscriptions
	var totalSubscriptions int64
	config.DB.Model(&models.Subscription{}).Where("is_active = ?", true).Count(&totalSubscriptions)

	c.JSON(http.StatusOK, gin.H{
		"date":                 today,
		"today_revenue":        totalRevenue,
		"today_payments":       len(payments),
		"active_subscriptions": totalSubscriptions,
	})
}
