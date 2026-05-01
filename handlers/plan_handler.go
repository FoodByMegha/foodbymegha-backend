package handlers

import (
	"net/http"
	"time"

	"github.com/FoodByMegha/foodbymegha-backend/config"
	"github.com/FoodByMegha/foodbymegha-backend/models"
	"github.com/gin-gonic/gin"
)

// GET /plans — sabko dikhao available plans
func GetPlans(c *gin.Context) {
	var plans []models.Plan
	config.DB.Where("is_active = ?", true).Find(&plans)
	c.JSON(http.StatusOK, gin.H{
		"plans": plans,
	})
}

// GET /my-plan — user apna plan dekhe
func GetMyPlan(c *gin.Context) {
	userID, _ := c.Get("userID")

	var subscription models.Subscription
	result := config.DB.Preload("Plan").Where("user_id = ? AND is_active = ?", userID, true).First(&subscription)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Koi active plan nahi hai"})
		return
	}

	// ✅ Auto-expire check
	if time.Now().After(subscription.EndDate) {
		config.DB.Model(&subscription).Update("is_active", false)
		c.JSON(http.StatusNotFound, gin.H{"error": "Tera plan khatam ho gaya! Naya plan lo 🍱"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"my_plan": subscription,
	})
}
