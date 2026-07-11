package handler

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go-perjalanan-dinas/config"
	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/src/service"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func isFinanceDepartment(departmentID uint) bool {
	if departmentID == 0 {
		return false
	}
	var dept models.Department
	if err := config.DB.First(&dept, departmentID).Error; err != nil {
		return false
	}
	code := strings.ToUpper(dept.Code)
	name := strings.ToLower(dept.Name)
	return code == "FIN" || strings.Contains(name, "finance") || strings.Contains(name, "keuangan")
}

type ReimbursementHandler struct {
	reimbursementService service.ReimbursementService
}

func NewReimbursementHandler(reimbursementService service.ReimbursementService) *ReimbursementHandler {
	return &ReimbursementHandler{reimbursementService: reimbursementService}
}

// SubmitClaim godoc
// @Summary Submit a travel expense claim
// @Tags claims
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param trip_id formData int true "Trip ID"
// @Param title formData string true "Claim title"
// @Param description formData string true "Claim description"
// @Param amount formData number true "Claim amount"
// @Param transaction_date formData string true "Transaction date (YYYY-MM-DD)"
// @Param files formData file false "Proof files"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /claims [post]
func (h *ReimbursementHandler) SubmitClaim(c *gin.Context) {
	var input dto.SubmitClaimInput
	if err := c.ShouldBind(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var files []*multipart.FileHeader
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		form, err := c.MultipartForm()
		if err == nil && form != nil {
			files = form.File["files"]
		}
	}

	userID, exists := c.Get("userID")
	if !exists {
		writeError(c, http.StatusUnauthorized, models.ErrUnauthorizedAccess.Error())
		return
	}

	claim, err := h.reimbursementService.SubmitClaimWithFiles(userID.(uint), input, files)
	if err != nil {
		slog.Error("Gagal submit klaim", "user_id", userID, "trip_id", input.TripID, "error", err)
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("Klaim berhasil diajukan", "claim_id", claim.ID, "user_id", userID, "trip_id", claim.TripID, "amount", claim.Amount)
	writeSuccess(c, http.StatusCreated, "Pengajuan klaim berhasil dibuat", claim)
}

// ReviewClaim godoc
// @Summary Review a travel expense claim (Admin Finance only)
// @Tags claims
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Claim ID"
// @Param request body dto.ReviewClaimInput true "Review payload"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /claims/{id}/review [patch]
func (h *ReimbursementHandler) ReviewClaim(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, models.ErrTripIDInvalid.Error())
		return
	}

	var input dto.ReviewClaimInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	role, _ := c.Get("role")
	departmentID, _ := c.Get("departmentID")
	var deptID uint
	if departmentID != nil {
		deptID = departmentID.(uint)
	}

	if role == "MANAGER" {
		if !isFinanceDepartment(deptID) {
			writeError(c, http.StatusForbidden, "Akses ditolak: Hanya manager departemen keuangan yang dapat menyetujui klaim")
			return
		}
	}

	claim, err := h.reimbursementService.ReviewClaim(uint(id), input)
	if err != nil {
		slog.Error("Gagal review klaim", "claim_id", id, "status", input.Status, "error", err)
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	message := "Klaim berhasil disetujui"
	if input.Status == "REJECTED" {
		message = "Klaim berhasil ditolak"
	}

	slog.Info("Review klaim berhasil", "claim_id", claim.ID, "status", input.Status)
	writeSuccess(c, http.StatusOK, message, claim)
}

// GetClaimsByTripID godoc
// @Summary Get claims for a trip
// @Tags claims
// @Produce json
// @Security BearerAuth
// @Param trip_id path int true "Trip ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /claims/trip/{trip_id} [get]
func (h *ReimbursementHandler) GetClaimsByTripID(c *gin.Context) {
	tripIDParam := c.Param("trip_id")
	tripID, err := strconv.ParseUint(tripIDParam, 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, models.ErrTripIDInvalid.Error())
		return
	}

	role, _ := c.Get("role")
	departmentID, _ := c.Get("departmentID")
	var deptID uint
	if departmentID != nil {
		deptID = departmentID.(uint)
	}

	if role == "MANAGER" {
		if !isFinanceDepartment(deptID) {
			var trip models.BusinessTrip
			if err := config.DB.Preload("User").First(&trip, tripID).Error; err != nil {
				writeError(c, http.StatusNotFound, "data perjalanan dinas tidak ditemukan")
				return
			}
			if trip.User.DepartmentID == nil || *trip.User.DepartmentID != deptID {
				writeError(c, http.StatusForbidden, "Akses ditolak: Anda hanya dapat melihat data klaim dari departemen Anda")
				return
			}
		}
	}

	claims, err := h.reimbursementService.GetClaimsByTripID(uint(tripID))
	if err != nil {
		writeError(c, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(c, http.StatusOK, "Daftar klaim berhasil diambil", claims)
}

func (h *ReimbursementHandler) DownloadClaimAttachment(c *gin.Context) {
	claimIDParam := c.Param("id")
	claimID, err := strconv.ParseUint(claimIDParam, 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, "ID klaim tidak valid")
		return
	}

	claim, err := h.reimbursementService.GetClaimByID(uint(claimID))
	if err != nil {
		writeError(c, http.StatusNotFound, "data klaim tidak ditemukan")
		return
	}

	role, _ := c.Get("role")
	departmentID, _ := c.Get("departmentID")
	var deptID uint
	if departmentID != nil {
		deptID = departmentID.(uint)
	}

	if role == "MANAGER" {
		if !isFinanceDepartment(deptID) {
			var trip models.BusinessTrip
			if err := config.DB.Preload("User").First(&trip, claim.TripID).Error; err != nil {
				writeError(c, http.StatusNotFound, "data perjalanan dinas tidak ditemukan")
				return
			}
			if trip.User.DepartmentID == nil || *trip.User.DepartmentID != deptID {
				writeError(c, http.StatusForbidden, "Akses ditolak: Anda hanya dapat melihat data klaim dari departemen Anda")
				return
			}
		}
	}

	filename := c.Param("filename")
	attachmentPaths := strings.Split(claim.AttachmentPaths, ",")
	matched := ""
	for _, path := range attachmentPaths {
		if path != "" && strings.HasSuffix(path, filename) {
			matched = path
			break
		}
	}

	if matched == "" {
		// Fallback to local scan if it wasn't tracked by the new path scheme (e.g. old local upload)
		claimDir := filepath.Join("uploads", "claims", fmt.Sprintf("claim-%d", claim.ID))
		if _, err := os.Stat(claimDir); err == nil {
			entries, err := os.ReadDir(claimDir)
			if err == nil {
				for _, entry := range entries {
					if entry.IsDir() {
						continue
					}
					if entry.Name() == filename || strings.HasSuffix(entry.Name(), filename) {
						matched = filepath.Join(claimDir, entry.Name())
						break
					}
				}
			}
		}
	}

	if matched == "" {
		writeError(c, http.StatusNotFound, "file bukti klaim tidak ditemukan")
		return
	}

	if config.MinioClient != nil && !strings.Contains(matched, "uploads/") {
		key := config.GetMinioKey(matched, "claims", claim.ID)
		object, err := config.MinioClient.GetObject(context.Background(), config.MinioBucket, key, minio.GetObjectOptions{})
		if err != nil {
			writeError(c, http.StatusInternalServerError, "Gagal mengambil file dari MinIO: "+err.Error())
			return
		}
		defer object.Close()

		info, err := object.Stat()
		if err != nil {
			writeError(c, http.StatusNotFound, "file bukti klaim tidak ditemukan di MinIO")
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(matched)))
		c.Header("Content-Type", info.ContentType)
		c.Header("Content-Length", strconv.FormatInt(info.Size, 10))
		c.DataFromReader(http.StatusOK, info.Size, info.ContentType, object, nil)
	} else {
		c.FileAttachment(matched, filepath.Base(matched))
	}
}

func (h *ReimbursementHandler) GetClaimByID(c *gin.Context) {
	claimIDParam := c.Param("id")
	claimID, err := strconv.ParseUint(claimIDParam, 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, "ID klaim tidak valid")
		return
	}

	claim, err := h.reimbursementService.GetClaimByID(uint(claimID))
	if err != nil {
		writeError(c, http.StatusNotFound, "data klaim tidak ditemukan")
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	departmentID, _ := c.Get("departmentID")

	switch role {
	case "EMPLOYEE":
		var trip models.BusinessTrip
		if err := config.DB.First(&trip, claim.TripID).Error; err != nil || trip.UserID != userID.(uint) {
			writeError(c, http.StatusForbidden, "Akses ditolak: Anda tidak memiliki akses ke klaim ini")
			return
		}
	case "MANAGER":
		var deptID uint
		if departmentID != nil {
			deptID = departmentID.(uint)
		}
		if !isFinanceDepartment(deptID) {
			var trip models.BusinessTrip
			if err := config.DB.Preload("User").First(&trip, claim.TripID).Error; err != nil {
				writeError(c, http.StatusNotFound, "data perjalanan dinas tidak ditemukan")
				return
			}
			if trip.User.DepartmentID == nil || *trip.User.DepartmentID != deptID {
				writeError(c, http.StatusForbidden, "Akses ditolak: Anda hanya dapat melihat data klaim dari departemen Anda")
				return
			}
		}
	}

	writeSuccess(c, http.StatusOK, "Detail klaim berhasil diambil", claim)
}

func (h *ReimbursementHandler) UpdateClaim(c *gin.Context) {
	claimIDParam := c.Param("id")
	claimID, err := strconv.ParseUint(claimIDParam, 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, "ID klaim tidak valid")
		return
	}

	var input dto.SubmitClaimInput
	if err := c.ShouldBind(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var files []*multipart.FileHeader
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		form, err := c.MultipartForm()
		if err == nil && form != nil {
			files = form.File["files"]
		}
	}

	userID, exists := c.Get("userID")
	if !exists {
		writeError(c, http.StatusUnauthorized, models.ErrUnauthorizedAccess.Error())
		return
	}

	claim, err := h.reimbursementService.UpdateClaimWithFiles(userID.(uint), uint(claimID), input, files)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccess(c, http.StatusOK, "Pengajuan klaim berhasil diperbarui", claim)
}

func (h *ReimbursementHandler) DeleteClaim(c *gin.Context) {
	claimIDParam := c.Param("id")
	claimID, err := strconv.ParseUint(claimIDParam, 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, "ID klaim tidak valid")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		writeError(c, http.StatusUnauthorized, models.ErrUnauthorizedAccess.Error())
		return
	}

	err = h.reimbursementService.DeleteClaim(userID.(uint), uint(claimID))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessMessage(c, http.StatusOK, "Pengajuan klaim berhasil dihapus")
}
