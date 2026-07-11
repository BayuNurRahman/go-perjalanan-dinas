package service

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/src/repository"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Helper to locate project root and load .env during tests
func loadTestEnv() {
	if _, exists := os.LookupEnv("APP_ENV_LOADED"); exists {
		return
	}
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	dir := wd
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			_ = godotenv.Load(filepath.Join(dir, ".env"))
			os.Setenv("APP_ENV_LOADED", "true")
			return
		}
		dir = filepath.Dir(dir)
	}
}

func setupTestDB(t *testing.T) *gorm.DB {
	loadTestEnv()
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "secret45"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "Traveldb"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5434"
	}
	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbName + " port=" + port + " sslmode=" + sslMode
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// AutoMigrate all models
	err = db.AutoMigrate(&models.Department{}, &models.User{}, &models.BusinessTrip{}, &models.Role{}, &models.BlacklistedToken{}, &models.Reimbursement{})
	if err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	return db
}

func TestFeaturesIntegration(t *testing.T) {
	db := setupTestDB(t)

	// Run all tests inside a transaction so the database remains clean
	tx := db.Begin()
	defer tx.Rollback()

	// Clean up existing conflicting test users first
	tx.Unscoped().Where("email IN ?", []string{"employeeA@mail.com", "managerA@mail.com", "managerB@mail.com"}).Delete(&models.User{})

	// 1. Setup mock role and departments using FirstOrCreate
	var roleEmp models.Role
	tx.Where(models.Role{Name: "EMPLOYEE"}).FirstOrCreate(&roleEmp)
	var roleMgr models.Role
	tx.Where(models.Role{Name: "MANAGER"}).FirstOrCreate(&roleMgr)

	var deptA models.Department
	tx.Where(models.Department{Name: "Department A", Code: "DEPT-A"}).FirstOrCreate(&deptA)
	var deptB models.Department
	tx.Where(models.Department{Name: "Department B", Code: "DEPT-B"}).FirstOrCreate(&deptB)

	// Hash password
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// 2. Create users
	empUser := models.User{
		Name:         "Employee A",
		Email:        "employeeA@mail.com",
		Password:     string(hashedPass),
		Role:         "EMPLOYEE",
		RoleID:       &roleEmp.ID,
		DepartmentID: &deptA.ID,
	}
	if err := tx.Create(&empUser).Error; err != nil {
		t.Fatalf("Failed to create test employee: %v", err)
	}

	mgrUserA := models.User{
		Name:         "Manager A",
		Email:        "managerA@mail.com",
		Password:     string(hashedPass),
		Role:         "MANAGER",
		RoleID:       &roleMgr.ID,
		DepartmentID: &deptA.ID,
	}
	if err := tx.Create(&mgrUserA).Error; err != nil {
		t.Fatalf("Failed to create test manager A: %v", err)
	}

	mgrUserB := models.User{
		Name:         "Manager B",
		Email:        "managerB@mail.com",
		Password:     string(hashedPass),
		Role:         "MANAGER",
		RoleID:       &roleMgr.ID,
		DepartmentID: &deptB.ID,
	}
	if err := tx.Create(&mgrUserB).Error; err != nil {
		t.Fatalf("Failed to create test manager B: %v", err)
	}

	// Clean up uploaded files after test finishes
	defer func() {
		_ = os.RemoveAll(filepath.Join(projectRoot(), "uploads", "trips", "user-"+string(rune(empUser.ID))))
	}()

	// Initialize repositories and services on transaction tx
	userRepo := repository.NewUserRepository(tx)
	tripRepo := repository.NewTripRepository(tx)
	reimbursementRepo := repository.NewReimbursementRepository(tx)

	authSvc := NewAuthServiceWithBlacklist(userRepo, nil)
	tripSvc := NewTripService(tripRepo)
	claimSvc := NewReimbursementService(reimbursementRepo)

	// --- FEATURE 1: Authentication ---
	t.Run("Auth_Login_Success", func(t *testing.T) {
		res, err := authSvc.Login(dto.LoginInput{
			Email:    empUser.Email,
			Password: "password123",
		})
		if err != nil {
			t.Fatalf("expected login to succeed, got error: %v", err)
		}
		if res.Email != empUser.Email {
			t.Errorf("expected email %s, got %s", empUser.Email, res.Email)
		}
	})

	t.Run("Auth_Login_WrongPassword", func(t *testing.T) {
		_, err := authSvc.Login(dto.LoginInput{
			Email:    empUser.Email,
			Password: "wrongpassword",
		})
		if err == nil {
			t.Fatalf("expected login to fail for wrong password")
		}
	})

	t.Run("Auth_Login_NonExistentEmail", func(t *testing.T) {
		_, err := authSvc.Login(dto.LoginInput{
			Email:    "nonexistent@mail.com",
			Password: "password123",
		})
		if err == nil {
			t.Fatalf("expected login to fail for non-existent email")
		}
	})

	// --- FEATURE 2: Business Trip Upload & Validation ---
	var trip models.BusinessTrip

	t.Run("Trip_Create_Success_WithMockFile", func(t *testing.T) {
		// Mock PDF file header
		fileHeader := createTestMultipartFileHeader(t, "surat-tugas.pdf", "application/pdf", []byte("%PDF-1.4\n%fake pdf content"))
		files := []*multipart.FileHeader{fileHeader}

		input := dto.CreateTripInput{
			Destination: "Surabaya",
			StartDate:   "2026-08-01",
			EndDate:     "2026-08-05",
			Description: "Dinas luar kota",
			Initiator:   "Employee A",
			Nomor_Surat: "SURAT-001",
		}

		var err error
		trip, err = tripSvc.CreateTripWithFiles(input, empUser.ID, files)
		if err != nil {
			t.Fatalf("expected trip creation to succeed, got: %v", err)
		}
		if trip.Destination != "Surabaya" || trip.Status != "PENDING" {
			t.Errorf("unexpected trip details: %+v", trip)
		}
	})

	t.Run("Trip_Create_Negative_InvalidDateRange", func(t *testing.T) {
		input := dto.CreateTripInput{
			Destination: "Surabaya",
			StartDate:   "2026-08-05",
			EndDate:     "2026-08-01", // End date before start date
			Description: "Dinas luar kota",
			Initiator:   "Employee A",
			Nomor_Surat: "SURAT-002",
		}
		_, err := tripSvc.CreateTripWithFiles(input, empUser.ID, nil)
		if err != models.ErrTripDateRangeInvalid {
			t.Fatalf("expected date range invalid error, got: %v", err)
		}
	})

	t.Run("Trip_Create_Negative_InvalidMimeExtension", func(t *testing.T) {
		// Mocking text file disguised as pdf
		fileHeader := createTestMultipartFileHeader(t, "harmful_executable.exe", "application/x-msdownload", []byte("malicious content"))
		files := []*multipart.FileHeader{fileHeader}

		input := dto.CreateTripInput{
			Destination: "Surabaya",
			StartDate:   "2026-08-01",
			EndDate:     "2026-08-05",
			Description: "Dinas luar kota",
			Initiator:   "Employee A",
			Nomor_Surat: "SURAT-003",
		}
		_, err := tripSvc.CreateTripWithFiles(input, empUser.ID, files)
		if err == nil {
			t.Fatalf("expected file creation to fail for invalid extension")
		}
	})

	// --- FEATURE 3: Department-Restricted Manager Approval ---
	t.Run("Manager_Approval_Negative_CrossDepartment", func(t *testing.T) {
		// Manager B (Dept B) tries to approve Trip 1 (owned by Employee A in Dept A)
		_, err := tripSvc.UpdateStatus(trip.ID, "APPROVED", nil, 0, "MANAGER", deptB.ID)
		if err != models.ErrTripStatusAccessDenied {
			t.Fatalf("expected access denied error for cross department approval, got: %v", err)
		}
	})

	t.Run("Manager_Approval_Positive_SameDepartment", func(t *testing.T) {
		// Manager A (Dept A) approves Trip 1 (owned by Employee A in Dept A)
		approvedTrip, err := tripSvc.UpdateStatus(trip.ID, "APPROVED", nil, 0, "MANAGER", deptA.ID)
		if err != nil {
			t.Fatalf("expected manager approval to succeed, got error: %v", err)
		}
		if approvedTrip.Status != "APPROVED" {
			t.Errorf("expected trip status to be APPROVED, got: %s", approvedTrip.Status)
		}
		trip = approvedTrip
	})

	t.Run("Employee_StartDuty_Success", func(t *testing.T) {
		onDutyTrip, err := tripSvc.UpdateStatus(trip.ID, "ON_DUTY", nil, empUser.ID, "EMPLOYEE")
		if err != nil {
			t.Fatalf("expected employee start duty to succeed, got: %v", err)
		}
		if onDutyTrip.Status != "ON_DUTY" {
			t.Errorf("expected status to be ON_DUTY, got %s", onDutyTrip.Status)
		}
		trip = onDutyTrip
	})

	t.Run("Employee_Complete_Success", func(t *testing.T) {
		fileHeader := createTestMultipartFileHeader(t, "bukti-dinas.pdf", "application/pdf", []byte("%PDF-1.4\n%fake proof content"))
		files := []*multipart.FileHeader{fileHeader}
		completedTrip, err := tripSvc.UpdateStatus(trip.ID, "COMPLETED", files, empUser.ID, "EMPLOYEE")
		if err != nil {
			t.Fatalf("expected employee complete to succeed, got: %v", err)
		}
		if completedTrip.Status != "COMPLETED" {
			t.Errorf("expected status to be COMPLETED, got %s", completedTrip.Status)
		}
		if !strings.Contains(completedTrip.AttachmentPaths, "bukti-dinas.pdf") {
			t.Errorf("expected proof attachment to be saved under AttachmentPaths, got: %s", completedTrip.AttachmentPaths)
		}
		trip = completedTrip
	})

	// --- FEATURE 4: Reimbursement Claims ---
	var claim models.Reimbursement

	t.Run("Claim_Submit_Success_WithMockFile", func(t *testing.T) {
		fileHeader := createTestMultipartFileHeader(t, "receipt.pdf", "application/pdf", []byte("%PDF-1.4\n%fake pdf content"))
		files := []*multipart.FileHeader{fileHeader}

		// Advance trip to completed so claim is valid
		trip.Status = "COMPLETED"
		tx.Save(&trip)

		input := dto.SubmitClaimInput{
			TripID:          trip.ID,
			Title:           "Makan Malam",
			Description:     "Makan malam dengan klien",
			Amount:          150000,
			TransactionDate: "2026-08-02",
		}

		var err error
		claim, err = claimSvc.SubmitClaimWithFiles(empUser.ID, input, files)
		if err != nil {
			t.Fatalf("expected claim submission to succeed, got: %v", err)
		}
		if claim.Title != "Makan Malam" || claim.Status != "PENDING" {
			t.Errorf("unexpected claim details: %+v", claim)
		}
	})

	t.Run("Claim_Submit_Negative_InvalidAmount", func(t *testing.T) {
		input := dto.SubmitClaimInput{
			TripID:          trip.ID,
			Title:           "Makan Malam",
			Description:     "Makan malam dengan klien",
			Amount:          -50000, // Negative amount
			TransactionDate: "2026-08-02",
		}
		_, err := claimSvc.SubmitClaimWithFiles(empUser.ID, input, nil)
		if err == nil {
			t.Fatalf("expected claim submission to fail for negative amount")
		}
	})

	t.Run("Claim_Review_Approval", func(t *testing.T) {
		reviewedClaim, err := claimSvc.ReviewClaim(claim.ID, dto.ReviewClaimInput{
			Status: "APPROVED",
		})
		if err != nil {
			t.Fatalf("expected claim review to succeed, got error: %v", err)
		}
		if reviewedClaim.Status != "APPROVED" {
			t.Errorf("expected claim status to be APPROVED, got: %s", reviewedClaim.Status)
		}
	})
}

func createTestMultipartFileHeader(t *testing.T, filename, contentType string, content []byte) *multipart.FileHeader {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", "form-data; name=\"files\"; filename=\""+filename+"\"")
	header.Set("Content-Type", contentType)

	part, err := writer.CreatePart(header)
	if err != nil {
		t.Fatalf("unable to create multipart part: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("unable to write multipart content: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("unable to close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(5 << 20); err != nil {
		t.Fatalf("unable to parse multipart form: %v", err)
	}

	return req.MultipartForm.File["files"][0]
}
