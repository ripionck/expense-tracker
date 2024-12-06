package middlewares

import (
	"net/http"

	"expense-tracker/pkg/helpers.go"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization Header Provided"})
			c.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Set("name", claims.Name)
		c.Set("uid", claims.Uid)

		c.Next()
	}
}
