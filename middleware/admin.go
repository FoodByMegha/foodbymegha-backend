package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Sirf admin access kar sakta hai!",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
