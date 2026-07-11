package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddlewareReturnsStandardErrorForMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware(nil))
	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	if w.Body.String() == "" || !strings.Contains(w.Body.String(), `"success":false`) {
		t.Fatalf("expected standard error envelope, got %s", w.Body.String())
	}
}

func TestRoleBlockMiddlewareReturnsStandardErrorForForbiddenRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", "EMPLOYEE")
		c.Next()
	})
	r.Use(RoleBlockMiddleware("SUPER_ADMIN"))
	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, w.Code)
	}

	if w.Body.String() == "" || !strings.Contains(w.Body.String(), `"success":false`) {
		t.Fatalf("expected standard error envelope, got %s", w.Body.String())
	}
}
