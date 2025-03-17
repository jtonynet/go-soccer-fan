package middleware

import (
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/util"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := util.ExtractTokenID(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Set("userID", token)

		c.Next()
	}
}
