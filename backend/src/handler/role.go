package handler

import (
	"net/http"
	"strconv"

	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/src/service"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// CreateRole godoc
// @Summary Create a new role (Super Admin only)
// @Tags roles
// @Accept json
// @Produce json
// @Param request body dto.CreateRoleInput true "Role payload"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var input dto.CreateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	roleValue, exists := c.Get("role")
	if !exists {
		writeError(c, http.StatusUnauthorized, models.ErrUnauthorizedAccess.Error())
		return
	}

	userRole, _ := roleValue.(string)
	role, err := h.roleService.CreateRole(userRole, input)
	if err != nil {
		if err == models.ErrSuperAdminOnly {
			writeError(c, http.StatusForbidden, err.Error())
			return
		}
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccess(c, http.StatusCreated, "Role berhasil dibuat", role)
}

// GetRoles godoc
// @Summary Get all roles
// @Tags roles
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Security BearerAuth
// @Router /roles [get]
func (h *RoleHandler) GetRoles(c *gin.Context) {
	roles, err := h.roleService.GetRoles()
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Gagal mengambil daftar role")
		return
	}

	writeSuccess(c, http.StatusOK, "Daftar role berhasil diambil", roles)
}

// GetRoleByID godoc
// @Summary Get role by ID
// @Tags roles
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /roles/{id} [get]
func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, models.ErrRoleIDInvalid.Error())
		return
	}

	role, err := h.roleService.GetRoleByID(uint(id))
	if err != nil {
		writeError(c, http.StatusNotFound, models.ErrRoleNotFound.Error())
		return
	}

	writeSuccess(c, http.StatusOK, "Role berhasil diambil", role)
}

// UpdateRole godoc
// @Summary Update a role (Super Admin only)
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param request body dto.CreateRoleInput true "Role payload"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	var input dto.CreateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, models.ErrRoleIDInvalid.Error())
		return
	}

	role, err := h.roleService.UpdateRole(uint(id), input)
	if err != nil {
		writeError(c, http.StatusNotFound, err.Error())
		return
	}

	writeSuccess(c, http.StatusOK, "Role berhasil diperbarui", role)
}

// DeleteRole godoc
// @Summary Delete a role (Super Admin only)
// @Tags roles
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, models.ErrRoleIDInvalid.Error())
		return
	}

	err = h.roleService.DeleteRole(uint(id))
	if err != nil {
		writeError(c, http.StatusNotFound, err.Error())
		return
	}

	writeSuccessMessage(c, http.StatusOK, "Role berhasil dihapus")
}
