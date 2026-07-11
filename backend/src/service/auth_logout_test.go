package service

import (
	"testing"
	"time"

	"go-perjalanan-dinas/dto"
)

type fakeBlacklistedTokenRepo struct {
	blacklisted []string
}

func (f *fakeBlacklistedTokenRepo) Add(token string, expiresAt time.Time) error {
	f.blacklisted = append(f.blacklisted, token)
	return nil
}

func (f *fakeBlacklistedTokenRepo) IsBlacklisted(token string) (bool, error) {
	for _, item := range f.blacklisted {
		if item == token {
			return true, nil
		}
	}
	return false, nil
}

func TestAuthServiceLogoutRevokesToken(t *testing.T) {
	blacklistRepo := &fakeBlacklistedTokenRepo{}
	svc := &authService{blacklistedTokenRepository: blacklistRepo}

	err := svc.Logout("test-token")
	if err != nil {
		t.Fatalf("expected logout to succeed, got %v", err)
	}

	blacklisted, err := blacklistRepo.IsBlacklisted("test-token")
	if err != nil {
		t.Fatalf("expected blacklist check to succeed, got %v", err)
	}
	if !blacklisted {
		t.Fatal("expected token to be blacklisted after logout")
	}
}

func TestAuthServiceRegisterStillUsesUserRepository(t *testing.T) {
	userRepo := &fakeUserRepo{}
	blacklistRepo := &fakeBlacklistedTokenRepo{}
	svc := &authService{userRepository: userRepo, blacklistedTokenRepository: blacklistRepo}

	user, err := svc.Register(dto.RegisterInput{Name: "Admin", Email: "admin@example.com", Password: "password123", RoleID: 1, DepartmentID: 1})
	if err != nil {
		t.Fatalf("expected register to succeed, got %v", err)
	}
	if user.Email != "admin@example.com" {
		t.Fatalf("expected registered user email to be preserved, got %s", user.Email)
	}
}
