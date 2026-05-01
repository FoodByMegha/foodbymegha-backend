package handlers

import (
	"net/http"

	"github.com/FoodByMegha/foodbymegha-backend/config"
	"github.com/FoodByMegha/foodbymegha-backend/models"
	"github.com/FoodByMegha/foodbymegha-backend/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register — naya customer account banana
func Register(c *gin.Context) {
	// Step 1: Request se data lo
	var input struct {
		Name     string `json:"name" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Address  string `json:"address"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Sab fields bharo — naam, phone, email, password",
		})
		return
	}

	// Step 2: Password hash karo
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password), bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Password process nahi ho paya",
		})
		return
	}

	// Step 3: User banao
	user := models.User{
		Name:     input.Name,
		Phone:    input.Phone,
		Email:    input.Email,
		Password: string(hashedPassword),
		Address:  input.Address,
		Role:     "customer",
	}

	// Step 4: Database mein save karo
	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Account nahi ban paya — phone ya email already registered hai",
		})
		return
	}

	// Step 5: Success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Account ban gaya! 🎉",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"phone": user.Phone,
			"role":  user.Role,
		},
	})
}

// Login — existing customer login kare
func Login(c *gin.Context) {
	// Step 1: Input lo
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email aur password dono chahiye",
		})
		return
	}

	// Step 2: User dhundo DB mein
	var user models.User
	result := config.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email ya password galat hai",
		})
		return
	}

	// Step 3: Password check karo
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email ya password galat hai",
		})
		return
	}

	// Step 4: JWT token banao
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Token nahi ban paya",
		})
		return
	}

	// Step 5: Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful! 🎉",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Profile — logged in customer ka profile
func GetProfile(c *gin.Context) {
	// Middleware se UserID lo
	userID := c.MustGet("userID").(uint)

	// DB se user dhundo
	var user models.User
	result := config.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User nahi mila",
		})
		return
	}

	// Profile return karo
	c.JSON(http.StatusOK, gin.H{
		"id":      user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"phone":   user.Phone,
		"address": user.Address,
		"role":    user.Role,
	})
}
