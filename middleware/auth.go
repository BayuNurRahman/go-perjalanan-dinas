package middleware

import (
	"go-perjalanan-dinas/config"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token tidak ditemukan"})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.ParseWithClaims(tokenString, &config.JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
			return config.JWTSecretKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token tidak valid"})
			return
		}

		claims := token.Claims.(*config.JWTClaim)
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}