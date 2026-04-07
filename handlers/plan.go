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

// POST /subscribe — user plan leta hai
func Subscribe(c *gin.Context) {
	// JWT se user ID nikalo
	userID, _ := c.Get("userID")

	// Request body se plan_id lo
	var input struct {
		PlanID uint `json:"plan_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plan ID dena zaroori hai"})
		return
	}

	// Plan exist karta hai?
	var plan models.Plan
	if err := config.DB.First(&plan, input.PlanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan nahi mila"})
		return
	}

	// Pehle se subscription hai?
	var existing models.Subscription
	config.DB.Where("user_id = ? AND is_active = ?", userID, true).First(&existing)
	if existing.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tumhara plan pehle se active hai"})
		return
	}

	// Naya subscription banao
	now := time.Now()
	subscription := models.Subscription{
		UserID:    userID.(uint),
		PlanID:    plan.ID,
		StartDate: now,
		EndDate:   now.AddDate(0, 0, plan.DurationDays),
		IsActive:  true,
	}

	config.DB.Create(&subscription)
	c.JSON(http.StatusOK, gin.H{
		"message":      "Plan le liya! Tiffin aayega 🍱",
		"subscription": subscription,
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

	c.JSON(http.StatusOK, gin.H{
		"my_plan": subscription,
	})
}
