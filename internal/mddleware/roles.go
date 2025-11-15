package mddleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.ToUpper(strings.TrimSpace(c.GetHeader("token")))
		if token != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "admin access required",
				},
			})
			c.Abort()
			return
		}
		c.Set("role", token)
		c.Next()
	}
}

func RequireAdminOrUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimSpace(c.GetHeader("token"))
		token = strings.ToUpper(token)
		if token != "ADMIN" && token != "USER" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "token required: ADMIN or USER",
				},
			})
			c.Abort()
			return
		}
		c.Set("role", token)
		c.Next()
	}
}
