package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/FoodByMegha/foodbymegha-backend/config"
	"github.com/FoodByMegha/foodbymegha-backend/models"
	"github.com/FoodByMegha/foodbymegha-backend/utils"
	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

// POST /payment — Razorpay order banao
func CreatePayment(c *gin.Context) {
	userID, _ := c.Get("userID")

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

	// Pehle se active subscription hai?
	var existing models.Subscription
	result := config.DB.Where("user_id = ? AND is_active = ?", userID, true).First(&existing)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tumhara plan pehle se active hai!"})
		return
	}

	// Razorpay client banao
	client := razorpay.NewClient(
		os.Getenv("RAZORPAY_KEY_ID"),
		os.Getenv("RAZORPAY_KEY_SECRET"),
	)

	// Amount paisa mein — ₹999 = 99900 paisa
	amountInPaise := int64(plan.Price * 100)

	data := map[string]interface{}{
		"amount":   amountInPaise,
		"currency": "INR",
		"receipt":  fmt.Sprintf("receipt_plan_%d_user_%d", plan.ID, userID),
	}

	rzpOrder, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Razorpay order nahi bana"})
		return
	}

	// Payment record DB mein save karo
	payment := models.Payment{
		UserID:          userID.(uint),
		PlanID:          plan.ID,
		RazorpayOrderID: rzpOrder["id"].(string),
		Amount:          plan.Price,
		Currency:        "INR",
		Status:          "created",
	}
	config.DB.Create(&payment)

	c.JSON(http.StatusOK, gin.H{
		"razorpay_order_id": rzpOrder["id"],
		"amount":            amountInPaise,
		"currency":          "INR",
		"key_id":            os.Getenv("RAZORPAY_KEY_ID"),
		"plan_name":         plan.Name,
	})
}

// POST /payment/verify — Payment verify karo aur subscription activate karo
func VerifyPayment(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input struct {
		RazorpayOrderID   string `json:"razorpay_order_id" binding:"required"`
		RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
		RazorpaySignature string `json:"razorpay_signature" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sab fields bharo"})
		return
	}

	// Signature verify karo
	secret := os.Getenv("RAZORPAY_KEY_SECRET")
	data := input.RazorpayOrderID + "|" + input.RazorpayPaymentID
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	generatedSignature := hex.EncodeToString(h.Sum(nil))

	if generatedSignature != input.RazorpaySignature {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment verify nahi hui!"})
		return
	}

	// Payment DB mein update karo
	var payment models.Payment
	config.DB.Where("razorpay_order_id = ?", input.RazorpayOrderID).First(&payment)
	config.DB.Model(&payment).Updates(map[string]interface{}{
		"razorpay_payment_id": input.RazorpayPaymentID,
		"status":              "paid",
	})

	// ✅ Subscription activate karo
	var plan models.Plan
	config.DB.First(&plan, payment.PlanID)

	now := time.Now()
	subscription := models.Subscription{
		UserID:    userID.(uint),
		PlanID:    payment.PlanID,
		StartDate: now,
		EndDate:   now.AddDate(0, 0, plan.DurationDays),
		IsActive:  true,
	}
	config.DB.Create(&subscription)

	// ✅ Email bhejo
	var user models.User
	config.DB.First(&user, userID)
	go utils.PaymentSuccessEmail(user.Email, user.Name, plan.Name, payment.Amount)

	// ✅ Sirf ek response
	c.JSON(http.StatusOK, gin.H{
		"message": "Payment successful! Subscription active ho gaya! 🎉",
	})
}

// GET /payment/history
func GetPaymentHistory(c *gin.Context) {
	userID, _ := c.Get("userID")

	var payments []models.Payment
	config.DB.Preload("Plan").Where("user_id = ?", userID).Find(&payments)

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
	})
}
