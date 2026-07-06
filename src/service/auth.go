package service

import (
	"errors"
	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/src/repository"
	"golang.org/x/crypto/bcrypt"
)

// 1. Definisikan Interface
type AuthService interface {
	Register(input dto.RegisterInput) (models.User, error)
	Login(input dto.LoginInput) (models.User, error)
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo}
}

// 2. Logika Pendaftaran (Register)
func (s *authService) Register(input dto.RegisterInput) (models.User, error) {
    // Hash password menggunakan bcrypt (Cost 10 adalah standar aman)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return models.User{}, err
    }

    // Ambil role dari input, jika kosong pasang default EMPLOYEE
    role := input.Role
    if role == "" {
        role = "EMPLOYEE"
    }

    // Bentuk objek user baru
    user := models.User{
        Name:     input.Name,
        Email:    input.Email,
        Password: string(hashedPassword),
        Role:     role, // <--- Sekarang mengambil isi variabel 'role' yang dinamis
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
		return user, errors.New("email tidak terdaftar")
	}

	// Cocokkan password inputan dengan password hash di database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return user, errors.New("password salah")
	}

	return user, nil
}