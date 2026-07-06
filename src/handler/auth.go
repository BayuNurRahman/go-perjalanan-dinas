package handler

import (
	"net/http"

	"go-perjalanan-dinas/config"
	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/src/service"
	

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.RegisterInput
	
	// Tangkap body JSON dari client
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Lempar ke service
	user, err := h.authService.Register(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mendaftarkan akun"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success", 
		"message": "Registrasi berhasil", 
		"data": map[string]string{"name": user.Name, "email": user.Email},
	})
}

// Di dalam file src/handler/auth_handler.go

func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	user, err := h.authService.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Generate JWT Token
	token, err := config.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": map[string]interface{}{
			"name":  user.Name,
			"role":  user.Role,
			"token": token,
		},
	})
}