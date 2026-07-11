package middleware

import (
	"net/http"
	"strings"

	"go-perjalanan-dinas/config"
	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/src/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(blacklistRepo repository.BlacklistedTokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Success: false, Message: models.ErrTokenNotFound.Error()})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &config.JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
			return config.JWTSecretKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Success: false, Message: models.ErrTokenInvalid.Error()})
			return
		}

		if blacklistRepo != nil {
			blacklisted, err := blacklistRepo.IsBlacklisted(tokenString)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Success: false, Message: models.ErrTokenVerificationFailed.Error()})
				return
			}
			if blacklisted {
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Success: false, Message: models.ErrTokenRevoked.Error()})
				return
			}
		}

		claims := token.Claims.(*config.JWTClaim)
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("departmentID", claims.DepartmentID)
		c.Next()
	}
}
