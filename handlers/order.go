package handlers

import (
	"net/http"
	"time"

	"github.com/FoodByMegha/foodbymegha-backend/config"
	"github.com/FoodByMegha/foodbymegha-backend/models"
	"github.com/gin-gonic/gin"
)

// GET /orders — user ke saare orders dekho
func GetOrders(c *gin.Context) {
	userID, _ := c.Get("userID")

	var orders []models.Order
	config.DB.Preload("Subscription").Where("user_id = ?", userID).Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

// GET /track/:id — ek order track karo
func TrackOrder(c *gin.Context) {
	userID, _ := c.Get("userID")
	orderID := c.Param("id")

	var order models.Order
	result := config.DB.Where("id = ? AND user_id = ?", orderID, userID).First(&order)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order nahi mila"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}

// PUT /orders/:id/status — admin status update kare
func UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status dena zaroori hai"})
		return
	}

	// Valid status check
	validStatuses := map[string]bool{
		"pending":          true,
		"out_for_delivery": true,
		"delivered":        true,
	}
	if !validStatuses[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status sirf pending, out_for_delivery, ya delivered ho sakta hai"})
		return
	}

	var order models.Order
	result := config.DB.First(&order, orderID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order nahi mila"})
		return
	}

	config.DB.Model(&order).Update("status", input.Status)
	c.JSON(http.StatusOK, gin.H{
		"message": "Status update ho gaya! ✅",
		"order":   order,
	})
}

// POST /orders — naya order banao (auto daily)
func CreateOrder(c *gin.Context) {
	userID, _ := c.Get("userID")

	// Active subscription hai?
	var subscription models.Subscription
	result := config.DB.Where("user_id = ? AND is_active = ?", userID, true).First(&subscription)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pehle plan lo!"})
		return
	}

	var input struct {
		Address string `json:"address" binding:"required"`
		Notes   string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address dena zaroori hai"})
		return
	}

	order := models.Order{
		UserID:         userID.(uint),
		SubscriptionID: subscription.ID,
		DeliveryDate:   time.Now().Format("2006-01-02"),
		Status:         "pending",
		Address:        input.Address,
		Notes:          input.Notes,
	}

	config.DB.Create(&order)
	c.JSON(http.StatusOK, gin.H{
		"message": "Order place ho gaya! 🍱",
		"order":   order,
	})
}
