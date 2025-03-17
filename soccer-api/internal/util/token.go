package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var secretKey = os.Getenv("API_SECRET_KEY")

func GenerateToken(userID uint) (string, error) {

	token_lifespan, err := strconv.Atoi(os.Getenv("API_TOKEN_LIFESPAN"))

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ExtractTokenID extracts the user ID from the JWT token
func ExtractTokenID(tokenString string) (uint, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, fmt.Errorf("invalid token %s", tokenString)
}

// TokenCookie sets the JWT token in a cookie
func TokenCookie(c *gin.Context) {
	// Get the token from the request header
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	// Set the token in the cookie
	c.SetCookie("token", tokenString, 0, "/", "localhost", false, true)
}
