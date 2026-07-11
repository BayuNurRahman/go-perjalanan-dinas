package middleware

import (
	"net/http"

	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"

	"github.com/gin-gonic/gin"
)

func RoleBlockMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Success: false, Message: models.ErrRoleContextMissing.Error()})
			return
		}

		isAllowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{Success: false, Message: models.ErrRoleForbidden.Error()})
			return
		}

		c.Next()
	}
}
