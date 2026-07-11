package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"go-perjalanan-dinas/config"
	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/src/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register a new user (Super Admin only, or first user if no Super Admin exists)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterInput true "Registration payload"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.RegisterInput

	// Tangkap body JSON dari client
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Authorization is handled by RoleBlockMiddleware in routes.go
	// Service layer receives RoleID from DTO without any string conversion

	// Lempar ke service dengan RoleID dari DTO
	user, err := h.authService.Register(input)
	if err != nil {
		if err == models.ErrEmailAlreadyExists {
			slog.Warn("Registrasi gagal: email sudah ada", "email", input.Email)
			writeError(c, http.StatusBadRequest, models.ErrEmailAlreadyExists.Error())
			return
		}
		slog.Error("Registrasi gagal", "email", input.Email, "error", err)
		writeError(c, http.StatusInternalServerError, "Gagal mendaftarkan akun")
		return
	}

	slog.Info("Registrasi berhasil", "email", user.Email, "name", user.Name)
	writeSuccess(c, http.StatusCreated, "Registrasi berhasil", map[string]string{"name": user.Name, "email": user.Email})
}

// Login godoc
// @Summary Login an existing user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginInput true "Login payload"
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.Login(input)
	if err != nil {
		slog.Warn("Login gagal", "email", input.Email, "error", err)
		writeError(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Generate JWT Token
	var departmentID uint
	if user.DepartmentID != nil {
		departmentID = *user.DepartmentID
	}

	token, err := config.GenerateToken(user.ID, user.Email, user.Role, departmentID)
	if err != nil {
		slog.Error("Gagal generate token", "user_id", user.ID, "error", err)
		writeError(c, http.StatusInternalServerError, "Gagal generate token")
		return
	}

	slog.Info("Login berhasil", "user_id", user.ID, "email", user.Email, "role", user.Role)
	writeSuccess(c, http.StatusOK, "Login berhasil", map[string]interface{}{
		"name":       user.Name,
		"role":       user.Role,
		"token":      token,
		"department": user.Department,
	})
}

// Logout godoc
// @Summary Logout current session
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		writeError(c, http.StatusUnauthorized, models.ErrTokenNotFound.Error())
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if err := h.authService.Logout(tokenString); err != nil {
		slog.Error("Logout gagal", "error", err)
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	userID, _ := c.Get("userID")
	slog.Info("Logout berhasil", "user_id", userID)
	writeSuccessMessage(c, http.StatusOK, "Logout berhasil")
}
