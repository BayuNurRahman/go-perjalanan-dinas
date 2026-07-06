package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func RoleBlockMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Akses ditolak: Sesi tidak memiliki data peran",
			})
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
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Akses ditolak: Peran Anda tidak memiliki otoritas",
			})
			return
		}

		c.Next()
	}
}