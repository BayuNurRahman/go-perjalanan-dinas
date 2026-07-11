package handler

import (
	"net/http"
	"strconv"

	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/src/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUsers godoc
// @Summary Get all users (Super Admin only)
// @Tags users
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetUsers()
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Gagal mengambil daftar pengguna")
		return
	}

	writeSuccess(c, http.StatusOK, "Daftar pengguna berhasil diambil", users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, models.ErrUserIDInvalid.Error())
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		writeError(c, http.StatusNotFound, models.ErrUserNotFound.Error())
		return
	}

	// Map to response DTO for security
	response := dto.UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Role:         user.Role,
		DepartmentID: user.DepartmentID,
	}

	writeSuccess(c, http.StatusOK, "Pengguna berhasil diambil", response)
}

// UpdateUser godoc
// @Summary Update user (Super Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body dto.UpdateUserInput true "User update payload"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var input dto.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, models.ErrUserIDInvalid.Error())
		return
	}

	user, err := h.userService.UpdateUser(uint(id), input)
	if err != nil {
		writeError(c, http.StatusNotFound, err.Error())
		return
	}

	response := dto.UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Role:         user.Role,
		DepartmentID: user.DepartmentID,
	}

	writeSuccess(c, http.StatusOK, "Pengguna berhasil diperbarui", response)
}

// DeleteUser godoc
// @Summary Delete user (Super Admin only)
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, models.ErrUserIDInvalid.Error())
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		writeError(c, http.StatusNotFound, err.Error())
		return
	}

	writeSuccessMessage(c, http.StatusOK, "Pengguna berhasil dihapus")
}
