package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/FoodByMegha/foodbymegha-backend/config"
	"github.com/FoodByMegha/foodbymegha-backend/models"
	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

// POST /payment — Razorpay order banao
func CreatePayment(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input struct {
		OrderID uint `json:"order_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID dena zaroori hai"})
		return
	}

	// Order exist karta hai?
	var order models.Order
	if err := config.DB.First(&order, input.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order nahi mila"})
		return
	}

	// Razorpay client banao
	client := razorpay.NewClient(
		os.Getenv("RAZORPAY_KEY_ID"),
		os.Getenv("RAZORPAY_KEY_SECRET"),
	)

	// Razorpay pe order banao
	// Amount paisa mein hota hai — ₹3999 = 399900 paisa
	data := map[string]interface{}{
		"amount":   399900,
		"currency": "INR",
		"receipt":  "receipt_order_" + string(rune(input.OrderID)),
	}
	rzpOrder, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Razorpay order nahi bana"})
		return
	}

	// Payment record DB mein save karo
	payment := models.Payment{
		UserID:          userID.(uint),
		OrderID:         input.OrderID,
		RazorpayOrderID: rzpOrder["id"].(string),
		Amount:          3999.00,
		Currency:        "INR",
		Status:          "created",
	}
	config.DB.Create(&payment)

	c.JSON(http.StatusOK, gin.H{
		"message":           "Payment order bana! Ab pay karo 💳",
		"razorpay_order_id": rzpOrder["id"],
		"amount":            399900,
		"currency":          "INR",
		"key_id":            os.Getenv("RAZORPAY_KEY_ID"),
	})
}

// POST /payment/verify — Payment verify karo
func VerifyPayment(c *gin.Context) {
	var input struct {
		RazorpayOrderID   string `json:"razorpay_order_id" binding:"required"`
		RazorpayPaymentID string `json:"razorpay_payment_id" binding:"required"`
		RazorpaySignature string `json:"razorpay_signature" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sab fields bharo"})
		return
	}

	// Signature verify karo — ye security check hai
	secret := os.Getenv("RAZORPAY_KEY_SECRET")
	data := input.RazorpayOrderID + "|" + input.RazorpayPaymentID
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	generatedSignature := hex.EncodeToString(h.Sum(nil))

	if generatedSignature != input.RazorpaySignature {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment verify nahi hui — fraud ho sakta hai!"})
		return
	}

	// DB mein payment update karo
	var payment models.Payment
	config.DB.Where("razorpay_order_id = ?", input.RazorpayOrderID).First(&payment)
	config.DB.Model(&payment).Updates(map[string]interface{}{
		"razorpay_payment_id": input.RazorpayPaymentID,
		"status":              "paid",
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment successful! Tiffin pakka! 🎉",
	})
}

// GET /payment/history — payment history dekho
func GetPaymentHistory(c *gin.Context) {
	userID, _ := c.Get("userID")

	var payments []models.Payment
	config.DB.Where("user_id = ?", userID).Find(&payments)

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
	})
}
