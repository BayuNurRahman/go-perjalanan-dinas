package service

import (
	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/src/repository"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// 1. Definisikan Interface
type AuthService interface {
	Register(input dto.RegisterInput) (models.User, error)
	Login(input dto.LoginInput) (models.User, error)
	Logout(token string) error
	HasSuperAdmin() (bool, error)
}

type authService struct {
	userRepository             repository.UserRepository
	blacklistedTokenRepository repository.BlacklistedTokenRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepository: userRepo}
}

func NewAuthServiceWithBlacklist(userRepo repository.UserRepository, blacklistRepo repository.BlacklistedTokenRepository) AuthService {
	return &authService{userRepository: userRepo, blacklistedTokenRepository: blacklistRepo}
}

// 2. Logika Pendaftaran (Register) - Super Admin only
func (s *authService) Register(input dto.RegisterInput) (models.User, error) {
	// Check if any Super Admin exists
	hasSuperAdmin, err := s.HasSuperAdmin()
	if err != nil {
		return models.User{}, err
	}

	// If Super Admin exists, only Super Admin can register new users
	// Authorization check should be done via middleware based on JWT role
	if hasSuperAdmin {
		// Middleware (rbac.go) will enforce that only SUPER_ADMIN can call this endpoint
		// Service layer assumes authorization is already validated
	}

	_, err = s.userRepository.FindByEmail(input.Email)
	if err == nil {
		return models.User{}, models.ErrEmailAlreadyExists
	}

	// Hash password menggunakan bcrypt (Cost 10 adalah standar aman)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	// Bentuk objek user baru dengan RoleID dari DTO (tanpa default/fallback values)
	departmentID := input.DepartmentID
	roleID := input.RoleID

	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		Password:     string(hashedPassword),
		RoleID:       &roleID,
		DepartmentID: &departmentID,
	}

	// Lempar ke repository untuk disimpan ke database
	newUser, err := s.userRepository.Create(user)
	return newUser, err
}

// 3. Logika Masuk (Login)
func (s *authService) Login(input dto.LoginInput) (models.User, error) {
	// Cari user berdasarkan email
	user, err := s.userRepository.FindByEmail(input.Email)
	if err != nil {
		return user, models.ErrEmailNotRegistered
	}

	// Cocokkan password inputan dengan password hash di database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return user, models.ErrInvalidPassword
	}
	return user, nil
}

func (s *authService) Logout(token string) error {
	if token == "" {
		return models.ErrTokenRequired
	}
	if s.blacklistedTokenRepository == nil {
		return nil
	}
	return s.blacklistedTokenRepository.Add(token, time.Now().Add(24*time.Hour))
}

// HasSuperAdmin checks if any Super Admin user exists by verifying RoleID relationships
func (s *authService) HasSuperAdmin() (bool, error) {
	users, err := s.userRepository.FindAll()
	if err != nil {
		return false, err
	}

	for _, user := range users {
		if user.Role == "SUPER_ADMIN" {
			return true, nil
		}
	}
	return false, nil
}

// normalizeRole converts role string to standard uppercase format
// Used by other services for role validation
func normalizeRole(role string) string {
	return strings.ToUpper(strings.TrimSpace(strings.ReplaceAll(role, " ", "_")))
}
