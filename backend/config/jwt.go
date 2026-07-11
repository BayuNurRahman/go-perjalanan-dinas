package config

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Gunakan secret key yang kompleks di lingkungan produksi (sebaiknya ambil dari .env)
var JWTSecretKey = []byte(getJWTSecret())

func getJWTSecret() string {
	LoadEnv()
	return GetEnv("JWT_SECRET", "super_secret_key_ta_bayu")
}

type JWTClaim struct {
	UserID       uint   `json:"user_id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	DepartmentID uint   `json:"department_id,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, email string, role string, departmentID uint) (string, error) {
	claims := JWTClaim{
		UserID:       userID,
		Email:        email,
		Role:         role,
		DepartmentID: departmentID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token berlaku 24 jam
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecretKey)
}
