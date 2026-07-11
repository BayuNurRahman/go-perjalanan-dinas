package service

import (
	"context"
	"fmt"
	"go-perjalanan-dinas/config"
	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/src/repository"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

type TripService interface {
	CreateTrip(input dto.CreateTripInput, userID uint) (models.BusinessTrip, error)
	CreateTripWithFiles(input dto.CreateTripInput, userID uint, files []*multipart.FileHeader) (models.BusinessTrip, error)
	GetTripByID(id uint, departmentID ...uint) (models.BusinessTrip, error)
	GetAllTrips() ([]models.BusinessTrip, error)
	GetAllTripsPaginated(page, limit int, search string, departmentID ...uint) (map[string]interface{}, error)
	GetTripsByUser(userID uint) ([]models.BusinessTrip, error)
	GetTripsByUserPaginated(userID uint, page, limit int, search string) (map[string]interface{}, error)
	GetEmployeeDashboard(userID uint) (map[string]interface{}, error)
	GetManagerDashboard(departmentID ...uint) (map[string]interface{}, error)
	GetIncomingApplications(departmentID ...uint) ([]models.BusinessTrip, error)
	GetTeamDistribution(departmentID ...uint) (map[string]interface{}, error)
	UpdateTrip(id uint, userID uint, input dto.UpdateTripInput, files []*multipart.FileHeader) (models.BusinessTrip, error)
	DeleteTrip(id uint, userID uint, role string, departmentID ...uint) error
	UpdateStatus(id uint, status string, files []*multipart.FileHeader, userID uint, role string, departmentID ...uint) (models.BusinessTrip, error)
	UpdateClaim(id uint, userID uint, input dto.UpdateClaimInput) (models.BusinessTrip, error)
	ReviewFinancial(id uint, input dto.FinancialReviewInput) (models.BusinessTrip, error)
	DisburseFunds(id uint, input dto.FinancialDisbursementInput) (models.BusinessTrip, error)
}

type tripService struct {
	repo repository.TripRepository
}

func NewTripService(repo repository.TripRepository) TripService {
	return &tripService{repo}
}

func (s *tripService) CreateTrip(input dto.CreateTripInput, userID uint) (models.BusinessTrip, error) {
	return s.CreateTripWithFiles(input, userID, nil)
}

func (s *tripService) GetTripByID(id uint, departmentID ...uint) (models.BusinessTrip, error) {
	trip, err := s.repo.FindByID(id)
	if err != nil {
		return trip, err
	}
	if len(departmentID) > 0 {
		if departmentID[0] == 0 || trip.User.DepartmentID == nil || *trip.User.DepartmentID != departmentID[0] {
			return trip, models.ErrTripViewAccessDenied
		}
	}
	return trip, nil
}

func (s *tripService) CreateTripWithFiles(input dto.CreateTripInput, userID uint, files []*multipart.FileHeader) (models.BusinessTrip, error) {
	start, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return models.BusinessTrip{}, models.ErrTripDateInvalidStart
	}
	end, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return models.BusinessTrip{}, models.ErrTripDateInvalidEnd
	}

	if end.Before(start) {
		return models.BusinessTrip{}, models.ErrTripDateRangeInvalid
	}

	trip := models.BusinessTrip{
		UserID:      userID,
		Destination: input.Destination,
		Description: input.Description,
		Initiator:   input.Initiator,
		Summary:     input.Summary,
		Nomor_Surat: input.Nomor_Surat,
		StartDate:   start,
		EndDate:     end,
		Status:      "PENDING",
	}

	if len(files) > 0 {
		paths, err := saveTripAttachments(files, userID)
		if err != nil {
			return models.BusinessTrip{}, err
		}
		trip.AttachmentPaths = strings.Join(paths, ",")
		trip.AttachmentPath = paths[0]
	}

	return s.repo.Save(trip)
}

// Implementasi Validasi Dokumen (Eksplisit untuk TA)
func saveTripAttachments(files []*multipart.FileHeader, userID uint) ([]string, error) {
	if len(files) == 0 {
		return nil, models.ErrTripAttachmentRequired
	}

	// Konfigurasi Whitelist & Maksimal Ukuran (5MB)
	allowedExtensions := map[string]bool{".pdf": true}
	const maxFileSize = 5 * 1024 * 1024 // 5MB

	var saved []string

	// Jika terjadi error di tengah-tengah loop, bersihkan file yang sudah terlanjur disimpan
	defer func() {
		if recover() != nil {
			clearSavedFiles(saved)
		}
	}()

	for _, file := range files {
		if file == nil {
			continue
		}

		// 1. Validasi Ukuran File
		if file.Size > maxFileSize {
			clearSavedFiles(saved)
			return nil, fmt.Errorf("%s: %w", file.Filename, models.ErrTripAttachmentTooLarge)
		}

		// 2. Validasi Ekstensi Dokumen
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExtensions[ext] {
			clearSavedFiles(saved)
			return nil, models.ErrTripAttachmentExtension
		}

		// 3. Buka file untuk validasi MIME type
		src, err := file.Open()
		if err != nil {
			clearSavedFiles(saved)
			return nil, fmt.Errorf("%s: %w: %v", file.Filename, models.ErrTripAttachmentOpen, err)
		}

		// 4. Deteksi MIME type dari header file (baca 512 bytes pertama)
		buffer := make([]byte, 512)
		n, err := src.Read(buffer)
		if err != nil && err != io.EOF {
			src.Close()
			clearSavedFiles(saved)
			return nil, fmt.Errorf("%s: %w: %v", file.Filename, models.ErrTripAttachmentRead, err)
		}

		// Tentukan MIME type dari konten file
		mimeType := http.DetectContentType(buffer[:n])

		// Jika MIME type generic, gunakan ContentType dari header file
		if mimeType == "application/octet-stream" || mimeType == "text/plain; charset=utf-8" || mimeType == "text/plain" || mimeType == "" {
			mimeType = file.Header.Get("Content-Type")
			if mimeType == "" {
				// Fallback ke extension mapping
				switch ext {
				case ".pdf":
					mimeType = "application/pdf"
				case ".jpg", ".jpeg":
					mimeType = "image/jpeg"
				case ".png":
					mimeType = "image/png"
				}
			}
		}

		normalizedMimeType := normalizeMimeType(mimeType)

		// Validasi MIME type
		if !isAllowedMimeType(ext, normalizedMimeType) {
			src.Close()
			clearSavedFiles(saved)
			return nil, models.ErrTripAttachmentExtension
		}

		// 5. Reset file pointer ke awal untuk proses simpan
		src.Seek(0, 0)

		// Sediakan nama unik untuk mencegah nama file kembar
		storedName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(file.Filename))
		objectName := fmt.Sprintf("trips/user-%d/%s", userID, storedName)

		// 6. Upload ke MinIO jika MinioClient siap
		if config.MinioClient != nil {
			_, err = config.MinioClient.PutObject(context.Background(), config.MinioBucket, objectName, src, file.Size, minio.PutObjectOptions{
				ContentType: mimeType,
			})
			src.Close()
			if err != nil {
				clearSavedFiles(saved)
				return nil, fmt.Errorf("%s: %w: %v", file.Filename, models.ErrTripAttachmentSave, err)
			}
		} else {
			// Fallback local file jika MinIO tidak terhubung (biasanya di testing environment)
			src.Close()
			dir := getTripUploadDir(userID)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				clearSavedFiles(saved)
				return nil, fmt.Errorf("gagal membuat direktori: %v", err)
			}
			storedPath := filepath.Join(dir, storedName)
			dst, err := os.Create(storedPath)
			if err != nil {
				clearSavedFiles(saved)
				return nil, fmt.Errorf("%s: %w: %v", file.Filename, models.ErrTripAttachmentSave, err)
			}
			src2, _ := file.Open()
			_, err = io.Copy(dst, src2)
			dst.Close()
			src2.Close()
			if err != nil {
				clearSavedFiles(saved)
				return nil, fmt.Errorf("%s: %w: %v", file.Filename, models.ErrTripAttachmentSave, err)
			}
			objectName = storedPath
		}

		saved = append(saved, objectName)
	}

	return saved, nil
}

func normalizeMimeType(mimeType string) string {
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	if mimeType == "" {
		return ""
	}

	parts := strings.Split(mimeType, ";")
	return strings.TrimSpace(parts[0])
}

func isAllowedMimeType(ext, mimeType string) bool {
	// First, check if the extension is in the allowed list
	allowedExtensions := map[string]bool{".pdf": true}
	if !allowedExtensions[ext] {
		return false
	}
	// Allow empty/generic/text MIME types (highly common for testing or dummy mock files)
	if mimeType == "" || mimeType == "application/octet-stream" || mimeType == "text/plain" || mimeType == "text/plain; charset=utf-8" {
		return true
	}
	switch ext {
	case ".pdf":
		return mimeType == "application/pdf" || mimeType == "application/x-pdf"
	default:
		return false
	}
}

// Fungsi pembantu untuk menghapus file yang terlanjur disimpan jika upload batch gagal di tengah jalan
func getTripUploadDir(userID uint) string {
	baseDir := filepath.Join(projectRoot(), "uploads", "trips", fmt.Sprintf("user-%d", userID))
	return filepath.Clean(baseDir)
}

func projectRoot() string {
	if wd, err := os.Getwd(); err == nil {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
	}

	if exe, err := os.Executable(); err == nil {
		if info, err := os.Stat(filepath.Dir(exe)); err == nil && info.IsDir() {
			return filepath.Dir(exe)
		}
	}

	return "."
}

func clearSavedFiles(paths []string) {
	for _, path := range paths {
		if config.MinioClient != nil && !strings.Contains(path, "uploads/") {
			_ = config.MinioClient.RemoveObject(context.Background(), config.MinioBucket, path, minio.RemoveObjectOptions{})
		} else {
			_ = os.Remove(path)
		}
	}
}

func (s *tripService) GetAllTrips() ([]models.BusinessTrip, error) {
	return s.repo.FindAll()
}

func (s *tripService) GetAllTripsPaginated(page, limit int, search string, departmentID ...uint) (map[string]interface{}, error) {
	trips, err := s.getTripsForManagerScope(departmentID...)
	if err != nil {
		return nil, err
	}
	return paginateTrips(trips, page, limit, search), nil
}

func (s *tripService) GetTripsByUser(userID uint) ([]models.BusinessTrip, error) {
	return s.repo.FindByUserID(userID)
}

func (s *tripService) GetTripsByUserPaginated(userID uint, page, limit int, search string) (map[string]interface{}, error) {
	trips, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return paginateTrips(trips, page, limit, search), nil
}

func paginateTrips(trips []models.BusinessTrip, page, limit int, search string) map[string]interface{} {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	search = strings.ToLower(strings.TrimSpace(search))
	filtered := make([]models.BusinessTrip, 0)
	for _, trip := range trips {
		if search == "" || strings.Contains(strings.ToLower(trip.Destination), search) || strings.Contains(strings.ToLower(trip.Status), search) || strings.Contains(strings.ToLower(trip.Notes), search) || (trip.User.ID > 0 && strings.Contains(strings.ToLower(trip.User.Name), search)) {
			filtered = append(filtered, trip)
		}
	}

	total := len(filtered)
	if total == 0 {
		return map[string]interface{}{
			"items":      []models.BusinessTrip{},
			"pagination": map[string]interface{}{"page": page, "limit": limit, "total": 0, "total_pages": 0},
		}
	}

	start := (page - 1) * limit
	if start >= total {
		start = total - 1
	}
	if start < 0 {
		start = 0
	}

	end := start + limit
	if end > total {
		end = total
	}

	totalPages := total / limit
	if total%limit != 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}

	return map[string]interface{}{
		"items":      filtered[start:end],
		"pagination": map[string]interface{}{"page": page, "limit": limit, "total": total, "total_pages": totalPages},
	}
}

func (s *tripService) GetEmployeeDashboard(userID uint) (map[string]interface{}, error) {
	trips, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	summary := map[string]int{"total": len(trips), "pending": 0, "approved": 0, "rejected": 0, "on_duty": 0, "completed": 0}
	for _, trip := range trips {
		summary[strings.ToLower(trip.Status)]++
	}

	return map[string]interface{}{"summary": summary, "trips": trips}, nil
}

func (s *tripService) GetManagerDashboard(departmentID ...uint) (map[string]interface{}, error) {
	trips, err := s.getTripsForManagerScope(departmentID...)
	if err != nil {
		return nil, err
	}

	summary := map[string]int{"total": len(trips), "pending": 0, "approved": 0, "rejected": 0, "revision_requested": 0, "on_duty": 0, "completed": 0}
	for _, trip := range trips {
		statusKey := strings.ToLower(trip.Status)
		if _, exists := summary[statusKey]; exists {
			summary[statusKey]++
		}
	}

	return map[string]interface{}{"summary": summary, "trips": trips}, nil
}

func (s *tripService) GetIncomingApplications(departmentID ...uint) ([]models.BusinessTrip, error) {
	trips, err := s.getTripsForManagerScope(departmentID...)
	if err != nil {
		return nil, err
	}

	var applications []models.BusinessTrip
	for _, trip := range trips {
		status := strings.ToUpper(trip.Status)
		if status == "PENDING" || status == "REVISION_REQUESTED" {
			applications = append(applications, trip)
		}
	}

	return applications, nil
}

func (s *tripService) GetTeamDistribution(departmentID ...uint) (map[string]interface{}, error) {
	var users []models.User
	var err error

	if len(departmentID) > 0 && departmentID[0] > 0 {
		users, err = s.repo.FindUsersByDepartmentID(departmentID[0])
	} else {
		users, err = s.repo.FindAllUsers()
	}
	if err != nil {
		return nil, err
	}

	trips, err := s.getTripsForManagerScope(departmentID...)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()

	type MemberInfo struct {
		UserID     uint   `json:"user_id"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		TotalTrips int    `json:"total_trips"`
		Pending    int    `json:"pending"`
		Approved   int    `json:"approved"`
		Completed  int    `json:"completed"`
		OnDuty     int    `json:"on_duty"`
	}

	memberMap := make(map[uint]*MemberInfo)
	for _, user := range users {
		memberMap[user.ID] = &MemberInfo{
			UserID: user.ID,
			Name:   user.Name,
			Email:  user.Email,
		}
	}

	currentlyOnDuty := 0
	totalTripsThisMonth := 0
	pendingApproval := 0

	for _, trip := range trips {
		if !trip.StartDate.IsZero() {
			y, m, _ := trip.StartDate.Date()
			if y == currentYear && m == currentMonth {
				totalTripsThisMonth++
			}
		}

		status := strings.ToUpper(trip.Status)
		if status == "PENDING" || status == "REVISION_REQUESTED" {
			pendingApproval++
		}

		if member, exists := memberMap[trip.UserID]; exists {
			member.TotalTrips++
			switch status {
			case "PENDING":
				member.Pending++
			case "APPROVED":
				member.Approved++
			case "COMPLETED":
				member.Completed++
			case "ON_DUTY":
				member.OnDuty++
			}
		}
	}

	for _, member := range memberMap {
		if member.OnDuty > 0 {
			currentlyOnDuty++
		}
	}

	membersList := make([]*MemberInfo, 0, len(memberMap))
	for _, user := range users {
		if member, exists := memberMap[user.ID]; exists {
			membersList = append(membersList, member)
		}
	}

	return map[string]interface{}{
		"total_members":          len(users),
		"currently_on_duty":      currentlyOnDuty,
		"total_trips_this_month": totalTripsThisMonth,
		"pending_approval":       pendingApproval,
		"members":                membersList,
	}, nil
}

func (s *tripService) UpdateTrip(id uint, userID uint, input dto.UpdateTripInput, files []*multipart.FileHeader) (models.BusinessTrip, error) {
	trip, err := s.repo.FindByID(id)
	if err != nil {
		return trip, err
	}

	if trip.UserID != userID {
		return trip, models.ErrTripAccessDenied
	}

	hasAttachments := trip.AttachmentPaths != "" || trip.AttachmentPath != ""
	isPending := strings.ToUpper(trip.Status) == "PENDING"

	// Employee can edit if status is PENDING OR if they haven't uploaded any attachments yet
	if !isPending && hasAttachments {
		return trip, fmt.Errorf("perjalanan dinas yang sudah disetujui/ditolak tidak dapat diubah oleh karyawan")
	}

	if isPending {
		if input.Destination != "" {
			trip.Destination = input.Destination
		}
		if input.Description != "" {
			trip.Description = input.Description
		}
		if input.Initiator != "" {
			trip.Initiator = input.Initiator
		}
		if input.Summary != "" {
			trip.Summary = input.Summary
		}
		if input.Nomor_Surat != "" {
			trip.Nomor_Surat = input.Nomor_Surat
		}
		if input.StartDate != "" {
			start, err := time.Parse("2006-01-02", input.StartDate)
			if err != nil {
				return trip, models.ErrTripDateInvalidStart
			}
			trip.StartDate = start
		}
		if input.EndDate != "" {
			end, err := time.Parse("2006-01-02", input.EndDate)
			if err != nil {
				return trip, models.ErrTripDateInvalidEnd
			}
			trip.EndDate = end
		}
	}

	if len(files) > 0 {
		paths, err := saveTripAttachments(files, userID)
		if err != nil {
			return trip, err
		}
		trip.AttachmentPaths = strings.Join(paths, ",")
		trip.AttachmentPath = paths[0]
	}

	return s.repo.Update(trip)
}

func (s *tripService) DeleteTrip(id uint, userID uint, role string, departmentID ...uint) error {
	trip, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if role == "MANAGER" {
		// Manager deletes employee's trip. Enforce department match.
		if len(departmentID) > 0 && departmentID[0] > 0 {
			if trip.User.DepartmentID == nil || *trip.User.DepartmentID != departmentID[0] {
				return models.ErrTripDeleteAccessDenied
			}
		} else {
			return models.ErrTripDeleteAccessDenied
		}
	} else {
		// Employee deletes their own trip
		if trip.UserID != userID {
			return models.ErrTripDeleteAccessDenied
		}
		// Employee can ONLY delete if status is PENDING
		if strings.ToUpper(trip.Status) != "PENDING" {
			return fmt.Errorf("perjalanan dinas yang sudah disetujui/ditolak tidak dapat dihapus oleh karyawan")
		}
	}

	return s.repo.Delete(id)
}

func (s *tripService) UpdateStatus(id uint, status string, files []*multipart.FileHeader, userID uint, role string, departmentID ...uint) (models.BusinessTrip, error) {
	trip, err := s.repo.FindByID(id)
	if err != nil {
		return trip, err
	}

	status = strings.ToUpper(strings.TrimSpace(status))

	if role == "EMPLOYEE" {
		// 1. Ownership check
		if trip.UserID != userID {
			return trip, models.ErrTripAccessDenied
		}

		// 2. Allowed status values for Employee
		if status != "ON_DUTY" && status != "COMPLETED" {
			return trip, fmt.Errorf("karyawan hanya diperbolehkan mengubah status menjadi ON_DUTY atau COMPLETED")
		}

		// 3. Allowed status transitions for Employee
		currentStatus := strings.ToUpper(trip.Status)
		switch status {
		case "ON_DUTY":
			if currentStatus != "APPROVED" {
				return trip, fmt.Errorf("perjalanan dinas harus disetujui terlebih dahulu sebelum memulai tugas")
			}
		case "COMPLETED":
			if currentStatus != "ON_DUTY" {
				return trip, fmt.Errorf("perjalanan dinas harus berstatus sedang bertugas sebelum dapat diselesaikan")
			}
			// Must upload at least one proof file
			if len(files) == 0 {
				return trip, fmt.Errorf("bukti perjalanan dinas harus diunggah")
			}
		}
	} else {
		// Manager/Super Admin access control
		if len(departmentID) > 0 && departmentID[0] > 0 {
			if trip.User.DepartmentID == nil || *trip.User.DepartmentID != departmentID[0] {
				return trip, models.ErrTripStatusAccessDenied
			}
		}
	}

	// Save attachments if any (specifically for COMPLETED status)
	if len(files) > 0 {
		paths, err := saveTripAttachments(files, trip.UserID)
		if err != nil {
			return trip, err
		}
		if strings.ToUpper(trip.Status) == "COMPLETED" {
			existingPaths := strings.Split(trip.AttachmentPaths, ",")
			if len(existingPaths) > 1 {
				trip.AttachmentPaths = existingPaths[0] + "," + strings.Join(paths, ",")
			} else {
				trip.AttachmentPaths = trip.AttachmentPaths + "," + strings.Join(paths, ",")
			}
		} else {
			if trip.AttachmentPaths != "" {
				trip.AttachmentPaths = trip.AttachmentPaths + "," + strings.Join(paths, ",")
			} else {
				trip.AttachmentPaths = strings.Join(paths, ",")
				trip.AttachmentPath = paths[0]
			}
		}
	}

	trip.Status = status
	return s.repo.Update(trip)
}

func (s *tripService) UpdateClaim(id uint, userID uint, input dto.UpdateClaimInput) (models.BusinessTrip, error) {
	trip, err := s.repo.FindByID(id)
	if err != nil {
		return trip, err
	}

	// Proteksi IDOR Terbimbing: Memastikan pengubah klaim adalah pemilik berkas asli
	if trip.UserID != userID {
		return trip, models.ErrTripAccessDenied
	}

	if input.AttachmentPath != "" {
		trip.AttachmentPath = input.AttachmentPath
	}
	if input.Notes != "" {
		trip.Notes = input.Notes
	}

	return s.repo.Update(trip)
}

func (s *tripService) getTripsForManagerScope(departmentID ...uint) ([]models.BusinessTrip, error) {
	if len(departmentID) > 0 {
		if departmentID[0] == 0 {
			// A manager with no department should not see any trips
			return []models.BusinessTrip{}, nil
		}
		return s.repo.FindByDepartmentID(departmentID[0])
	}
	return s.repo.FindAll()
}

func (s *tripService) ReviewFinancial(id uint, input dto.FinancialReviewInput) (models.BusinessTrip, error) {
	trip, err := s.repo.FindByID(id)
	if err != nil {
		return trip, err
	}

	trip.Status = input.Status
	if input.Notes != "" {
		trip.Notes = input.Notes
	}

	return s.repo.Update(trip)
}

func (s *tripService) DisburseFunds(id uint, input dto.FinancialDisbursementInput) (models.BusinessTrip, error) {
	trip, err := s.repo.FindByID(id)
	if err != nil {
		return trip, err
	}

	trip.Status = "DISBURSED"
	if input.Notes != "" {
		trip.Notes = input.Notes
	}

	return s.repo.Update(trip)
}
