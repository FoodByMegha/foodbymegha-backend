package middleware

import (
	"net/http"
	"strings"

	"github.com/FoodByMegha/foodbymegha-backend/utils"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1: Header se token lo
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Pehle login karo! Token nahi mila",
			})
			c.Abort()
			return
		}

		// Step 2: "Bearer token123" se sirf token nikalo
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Step 3: Token verify karo
		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token invalid ya expire ho gaya — dobara login karo",
			})
			c.Abort()
			return
		}

		// Step 4: UserID aur Role save karo
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		// Step 5: Aage jaane do
		c.Next()
	}
}
