package handler

import (
	"context"
	"fmt"
	"go-perjalanan-dinas/config"
	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/src/service"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)



type TripHandler struct {
	tripService service.TripService
}

func NewTripHandler(tripService service.TripService) *TripHandler {
	return &TripHandler{tripService}
}

// CreateTrip godoc
// @Summary Create a travel request
// @Tags trips
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param destination formData string true "Destination"
// @Param start_date formData string true "Start date (YYYY-MM-DD)"
// @Param end_date formData string true "End date (YYYY-MM-DD)"
// @Param description formData string true "Description of trip"
// @Param initiator formData string true "Name of initiator"
// @Param summary formData string false "Summary of trip"
// @Param nomor_surat formData string true "Letter number"
// @Param files formData file false "Supporting files"
// @Success 201 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /trips [post]
func (h *TripHandler) CreateTrip(c *gin.Context) {
	var input dto.CreateTripInput

	// Gin otomatis mendeteksi Content-Type (JSON atau Form/Multipart)
	// dan mengikatnya sesuai tag `json` atau `form` di struct DTO kamu.
	if err := c.ShouldBind(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Ambil file secara spesifik jika request-nya adalah multipart
	var files []*multipart.FileHeader
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		form, err := c.MultipartForm()
		if err == nil && form != nil {
			files = form.File["files"]
		}

		// [Opsional] Validasi awal dokumen di tingkat Handler (Judul TA: Validasi Dokumen)
		// Kamu bisa cek ukuran file atau ekstensi di sini sebelum masuk service.
	}

	// Ambil userID dari JWT Middleware
	userID, exists := c.Get("userID")
	if !exists {
		writeError(c, http.StatusUnauthorized, models.ErrUnauthorizedAccess.Error())
		return
	}

	// Eksekusi ke layer Service/Usecase
	trip, err := h.tripService.CreateTripWithFiles(input, userID.(uint), files)
	if err != nil {
		slog.Error("Gagal membuat perjalanan dinas", "user_id", userID, "destination", input.Destination, "error", err)
		statusCode := http.StatusBadRequest
		if err == models.ErrUnauthorizedAccess {
			statusCode = http.StatusUnauthorized
		}
		writeError(c, statusCode, err.Error())
		return
	}

	slog.Info("Perjalanan dinas berhasil dibuat", "trip_id", trip.ID, "user_id", userID, "destination", trip.Destination)
	writeSuccess(c, http.StatusCreated, "Pengajuan perjalanan dinas berhasil dibuat", trip)
}

func (h *TripHandler) GetAllTrips(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	role, _ := c.Get("role")
	departmentID, _ := c.Get("departmentID")
	var deptID uint
	if departmentID != nil {
		deptID = departmentID.(uint)
	}

	var result map[string]interface{}
	var err error
	if role == "MANAGER" {
		result, err = h.tripService.GetAllTripsPaginated(page, limit, search, deptID)
	} else {
		result, err = h.tripService.GetAllTripsPaginated(page, limit, search)
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Gagal mengambil data monitoring")
		return
	}

	writeSuccess(c, http.StatusOK, "Data monitoring berhasil diambil", map[string]interface{}{
		"items":      result["items"],
		"pagination": result["pagination"],
	})
}

func (h *TripHandler) GetEmployeeDashboard(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	dashboard, err := h.tripService.GetEmployeeDashboard(userID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Gagal mengambil dashboard karyawan")
		return
	}

	writeSuccess(c, http.StatusOK, "Dashboard karyawan berhasil diambil", dashboard)
}

func (h *TripHandler) GetMyTrips(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	result, err := h.tripService.GetTripsByUserPaginated(userID, page, limit, search)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Gagal mengambil data perjalanan dinas Anda")
		return
	}

	writeSuccess(c, http.StatusOK, "Data perjalanan dinas Anda berhasil diambil", map[string]interface{}{
		"items":      result["items"],
		"pagination": result["pagination"],
	})
}

func (h *TripHandler) UpdateTrip(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, "ID perjalanan dinas tidak valid")
		return
	}

	var input dto.UpdateTripInput
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

	userID := c.MustGet("userID").(uint)
	trip, err := h.tripService.UpdateTrip(uint(id), userID, input, files)
	if err != nil {
		statusCode := http.StatusBadRequest
		message := err.Error()
		if strings.Contains(err.Error(), "akses ditolak") {
			statusCode = http.StatusForbidden
		}
		slog.Warn("Gagal memperbarui perjalanan dinas", "trip_id", id, "user_id", userID, "error", err)
		writeError(c, statusCode, message)
		return
	}

	slog.Info("Perjalanan dinas berhasil diperbarui", "trip_id", trip.ID, "user_id", userID)
	writeSuccess(c, http.StatusOK, "Perjalanan dinas berhasil diperbarui", trip)
}

func (h *TripHandler) DeleteTrip(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, "ID perjalanan dinas tidak valid")
		return
	}

	role, _ := c.Get("role")
	departmentID, _ := c.Get("departmentID")
	var deptID uint
	if departmentID != nil {
		deptID = departmentID.(uint)
	}

	userID := c.MustGet("userID").(uint)
	if role == "MANAGER" {
		err = h.tripService.DeleteTrip(uint(id), userID, role.(string), deptID)
	} else {
		err = h.tripService.DeleteTrip(uint(id), userID, role.(string))
	}
	if err != nil {
		statusCode := http.StatusBadRequest
		message := err.Error()
		if strings.Contains(err.Error(), "akses ditolak") {
			statusCode = http.StatusForbidden
		}
		slog.Warn("Gagal menghapus perjalanan dinas", "trip_id", id, "user_id", userID, "error", err)
		writeError(c, statusCode, message)
		return
	}

	slog.Info("Perjalanan dinas berhasil dihapus", "trip_id", id, "user_id", userID)
	writeSuccessMessage(c, http.StatusOK, "Perjalanan dinas berhasil dihapus")
}

func (h *TripHandler) UpdateClaim(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var input dto.UpdateClaimInput
	if err := c.ShouldBindJSON(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.MustGet("userID").(uint)
	trip, err := h.tripService.UpdateClaim(uint(id), userID, input)
	if err != nil {
		writeError(c, http.StatusForbidden, err.Error())
		return
	}

	writeSuccess(c, http.StatusOK, "Lampiran berhasil diperbarui", trip)
}

func (h *TripHandler) GetManagerDashboard(c *gin.Context) {
	role, _ := c.Get("role")
	departmentID, _ := c.Get("departmentID")
	var deptID uint
	if departmentID != nil {
		deptID = departmentID.(uint)
	}

	var dashboard map[string]interface{}
	var err error
	if role == "MANAGER" {
		dashboard, err = h.tripService.GetManagerDashboard(deptID)
	} else {
		dashboard, err = h.tripService.GetManagerDashboard()
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Gagal mengambil dashboard manajer")
		return
	}

	writeSuccess(c, http.StatusOK, "Dashboard manajer berhasil diambil", dashboard)
}

func (h *TripHandler) GetIncomingApplications(c *gin.Context) {
	role, _ := c.Get("role")
	departmentID, _ := c.Get("departmentID")
	var deptID uint
	if departmentID != nil {
		deptID = departmentID.(uint)
	}

	var applications []models.BusinessTrip
	var err error
	if role == "MANAGER" {
		applications, err = h.tripService.GetIncomingApplications(deptID)
	} else {
		applications, err = h.tripService.GetIncomingApplications()
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Gagal mengambil daftar aplikasi masuk")
		return
	}

	writeSuccess(c, http.StatusOK, "Daftar aplikasi masuk berhasil diambil", applications)
}

func (h *TripHandler) GetTeamDistribution(c *gin.Context) {
	role, _ := c.Get("role")
	departmentID, _ := c.Get("departmentID")
	var deptID uint
	if departmentID != nil {
		deptID = departmentID.(uint)
	}

	var distribution map[string]interface{}
	var err error
	if role == "MANAGER" {
		distribution, err = h.tripService.GetTeamDistribution(deptID)
	} else {
		distribution, err = h.tripService.GetTeamDistribution()
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "Gagal mengambil distribusi tim")
		return
	}

	writeSuccess(c, http.StatusOK, "Distribusi tim berhasil diambil", distribution)
}

func (h *TripHandler) ReviewFinancial(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var input dto.FinancialReviewInput
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
			writeError(c, http.StatusForbidden, "Akses ditolak: Hanya manager departemen keuangan yang dapat menyetujui/mereview klaim")
			return
		}
	}

	trip, err := h.tripService.ReviewFinancial(uint(id), input)
	if err != nil {
		slog.Error("Gagal review finansial", "trip_id", id, "status", input.Status, "error", err)
		writeError(c, http.StatusForbidden, err.Error())
		return
	}

	slog.Info("Review finansial berhasil", "trip_id", trip.ID, "status", input.Status)
	writeSuccess(c, http.StatusOK, "Review keuangan berhasil", trip)
}

func (h *TripHandler) DisburseFunds(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var input dto.FinancialDisbursementInput
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
			writeError(c, http.StatusForbidden, "Akses ditolak: Hanya manager departemen keuangan yang dapat mencairkan dana")
			return
		}
	}

	trip, err := h.tripService.DisburseFunds(uint(id), input)
	if err != nil {
		slog.Error("Gagal pencairan dana", "trip_id", id, "amount", input.Amount, "error", err)
		writeError(c, http.StatusForbidden, err.Error())
		return
	}

	slog.Info("Pencairan dana berhasil", "trip_id", trip.ID, "amount", input.Amount, "reference_id", input.ReferenceID)
	writeSuccess(c, http.StatusOK, "Pencairan dana berhasil", trip)
}

func (h *TripHandler) DownloadAttachment(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	role, _ := c.Get("role")
	departmentID, _ := c.Get("departmentID")
	var deptID uint
	if departmentID != nil {
		deptID = departmentID.(uint)
	}

	var trip models.BusinessTrip
	var getErr error
	if role == "MANAGER" {
		trip, getErr = h.tripService.GetTripByID(uint(id), deptID)
	} else {
		trip, getErr = h.tripService.GetTripByID(uint(id))
	}
	if getErr != nil {
		statusCode := http.StatusNotFound
		message := "data perjalanan dinas tidak ditemukan"
		if strings.Contains(getErr.Error(), "akses ditolak") {
			statusCode = http.StatusForbidden
			message = getErr.Error()
		}
		writeError(c, statusCode, message)
		return
	}

	filename := c.Param("filename")
	attachmentPaths := strings.Split(trip.AttachmentPaths, ",")
	matched := ""
	for _, path := range attachmentPaths {
		if path != "" && strings.HasSuffix(path, filename) {
			matched = path
		}
	}
	if matched == "" {
		writeError(c, http.StatusNotFound, "file lampiran tidak ditemukan")
		return
	}
	if config.MinioClient != nil && !strings.Contains(matched, "uploads/") {
		key := config.GetMinioKey(matched, "trips", trip.UserID)
		object, err := config.MinioClient.GetObject(context.Background(), config.MinioBucket, key, minio.GetObjectOptions{})
		if err != nil {
			writeError(c, http.StatusInternalServerError, "Gagal mengambil file dari MinIO: "+err.Error())
			return
		}
		defer object.Close()

		info, err := object.Stat()
		if err != nil {
			writeError(c, http.StatusNotFound, "file lampiran tidak ditemukan di MinIO")
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

func (h *TripHandler) UpdateStatus(c *gin.Context) {
	// 1. Ambil ID dari URL
	idParam := c.Param("id")

	// 2. Konversi string ke uint agar bisa diproses
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	// 3. Ambil status dan catatan dari request (bisa JSON atau Form/Multipart)
	var input struct {
		Status string `form:"status" json:"status" binding:"required"`
		Reason string `form:"reason" json:"reason"`
	}
	if err := c.ShouldBind(&input); err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Ambil file jika request-nya adalah multipart (unggah bukti dinas selesai)
	var files []*multipart.FileHeader
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		form, err := c.MultipartForm()
		if err == nil && form != nil {
			files = form.File["files"]
		}
	}

	// 5. Ambil data user dari JWT
	role, _ := c.Get("role")
	userID, _ := c.Get("userID")
	departmentID, _ := c.Get("departmentID")
	var deptID uint
	if departmentID != nil {
		deptID = departmentID.(uint)
	}

	var userUID uint
	if userID != nil {
		userUID = userID.(uint)
	}

	var trip models.BusinessTrip
	var roleStr string
	if role != nil {
		roleStr = role.(string)
	}

	if roleStr == "MANAGER" {
		trip, err = h.tripService.UpdateStatus(uint(id), input.Status, files, userUID, roleStr, deptID)
	} else {
		trip, err = h.tripService.UpdateStatus(uint(id), input.Status, files, userUID, roleStr)
	}

	if input.Reason != "" {
		trip.Notes = input.Reason
		updatedTrip, updateErr := h.tripService.UpdateClaim(uint(id), trip.UserID, dto.UpdateClaimInput{Notes: input.Reason})
		if updateErr != nil {
			err = updateErr
		} else {
			trip = updatedTrip
		}
	}

	if err != nil {
		statusCode := http.StatusBadRequest
		message := err.Error()
		if strings.Contains(err.Error(), "akses ditolak") {
			statusCode = http.StatusForbidden
		}
		slog.Warn("Gagal memperbarui status perjalanan dinas", "trip_id", id, "status", input.Status, "role", roleStr, "error", err)
		writeError(c, statusCode, message)
		return
	}

	slog.Info("Status perjalanan dinas berhasil diperbarui", "trip_id", trip.ID, "status", trip.Status, "role", roleStr)
	writeSuccess(c, http.StatusOK, "Status perjalanan dinas berhasil diperbarui", trip)
}

func (h *TripHandler) GetTripByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		writeError(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	role, _ := c.Get("role")
	departmentID, _ := c.Get("departmentID")
	var deptID uint
	if departmentID != nil {
		deptID = departmentID.(uint)
	}

	var trip models.BusinessTrip
	if role == "MANAGER" {
		trip, err = h.tripService.GetTripByID(uint(id), deptID)
	} else {
		trip, err = h.tripService.GetTripByID(uint(id))
	}

	if err != nil {
		statusCode := http.StatusNotFound
		message := "data perjalanan dinas tidak ditemukan"
		if strings.Contains(err.Error(), "akses ditolak") {
			statusCode = http.StatusForbidden
			message = err.Error()
		}
		writeError(c, statusCode, message)
		return
	}

	writeSuccess(c, http.StatusOK, "Detail perjalanan dinas berhasil diambil", trip)
}
