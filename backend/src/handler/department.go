package handler

import (
	"net/http"
	"strconv"

	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/src/service"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	departmentService service.DepartmentService
}

func NewDepartmentHandler(departmentService service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{departmentService: departmentService}
}

// CreateDepartment godoc
// @Summary Create a department (Super Admin only)
// @Tags departments
// @Accept json
// @Produce json
// @Param request body dto.CreateDepartmentInput true "Department payload"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /departments [post]
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var input dto.CreateDepartmentInput
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
	department, err := h.departmentService.CreateDepartment(userRole, input)
	if err != nil {
		writeError(c, http.StatusForbidden, err.Error())
		return
	}

	writeSuccess(c, http.StatusCreated, "Departemen berhasil dibuat", department)
}

// GetDepartments godoc
// @Summary Get all departments
// @Tags departments
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Security BearerAuth
// @Router /departments [get]
func (h *DepartmentHandler) GetDepartments(c *gin.Context) {
	departments, err := h.departmentService.GetDepartments()
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Gagal mengambil daftar departemen")
		return
	}

	writeSuccess(c, http.StatusOK, "Daftar departemen berhasil diambil", departments)
}

// GetDepartmentByID godoc
// @Summary Get department by ID
// @Tags departments
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /departments/{id} [get]
func (h *DepartmentHandler) GetDepartmentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, models.ErrDepartmentIDInvalid.Error())
		return
	}

	department, err := h.departmentService.GetDepartmentByID(uint(id))
	if err != nil {
		writeError(c, http.StatusNotFound, models.ErrDepartmentNotFound.Error())
		return
	}

	writeSuccess(c, http.StatusOK, "Departemen berhasil diambil", department)
}

// UpdateDepartment godoc
// @Summary Update a department (Super Admin only)
// @Tags departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param request body dto.CreateDepartmentInput true "Department payload"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /departments/{id} [put]
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	var input dto.CreateDepartmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, models.ErrDepartmentIDInvalid.Error())
		return
	}

	department, err := h.departmentService.UpdateDepartment(uint(id), input)
	if err != nil {
		writeError(c, http.StatusNotFound, err.Error())
		return
	}

	writeSuccess(c, http.StatusOK, "Departemen berhasil diperbarui", department)
}

// DeleteDepartment godoc
// @Summary Delete a department (Super Admin only)
// @Tags departments
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /departments/{id} [delete]
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, models.ErrDepartmentIDInvalid.Error())
		return
	}

	err = h.departmentService.DeleteDepartment(uint(id))
	if err != nil {
		writeError(c, http.StatusNotFound, err.Error())
		return
	}

	writeSuccessMessage(c, http.StatusOK, "Departemen berhasil dihapus")
}
