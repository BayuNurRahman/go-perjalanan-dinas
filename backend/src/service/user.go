package service

import (
	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/src/repository"
)

type UserService interface {
	GetUsers() ([]models.User, error)
	GetUserByID(id uint) (models.User, error)
	UpdateUser(id uint, input dto.UpdateUserInput) (models.User, error)
	DeleteUser(id uint) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

// GetUsers retrieves all users
func (s *userService) GetUsers() ([]models.User, error) {
	return s.userRepository.FindAll()
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id uint) (models.User, error) {
	if id == 0 {
		return models.User{}, models.ErrUserIDInvalid
	}
	return s.userRepository.FindByID(id)
}

// UpdateUser updates a user's profile
func (s *userService) UpdateUser(id uint, input dto.UpdateUserInput) (models.User, error) {
	if id == 0 {
		return models.User{}, models.ErrUserIDInvalid
	}

	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return models.User{}, models.ErrUserNotFound
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.RoleID != nil {
		user.RoleID = input.RoleID
	}
	if input.DepartmentID != nil {
		user.DepartmentID = input.DepartmentID
	}

	return s.userRepository.Update(user)
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(id uint) error {
	if id == 0 {
		return models.ErrUserIDInvalid
	}

	// Verify user exists before deleting
	_, err := s.userRepository.FindByID(id)
	if err != nil {
		return models.ErrUserNotFound
	}

	return s.userRepository.Delete(id)
}
