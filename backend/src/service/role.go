package service

import (
	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/src/repository"
)

type RoleService interface {
	CreateRole(userRole string, input dto.CreateRoleInput) (models.Role, error)
	GetRoles() ([]models.Role, error)
	GetRoleByID(id uint) (models.Role, error)
	GetRoleByName(name string) (models.Role, error)
	UpdateRole(id uint, input dto.CreateRoleInput) (models.Role, error)
	DeleteRole(id uint) error
}

type roleService struct {
	roleRepository repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{repo}
}

// CreateRole creates a new role (Super Admin only)
func (s *roleService) CreateRole(userRole string, input dto.CreateRoleInput) (models.Role, error) {
	if normalizeRole(userRole) != "SUPER_ADMIN" {
		return models.Role{}, models.ErrSuperAdminOnly
	}

	if input.Name == "" {
		return models.Role{}, models.ErrRoleNameRequired
	}

	// Check if role already exists
	normalizedName := normalizeRole(input.Name)
	existingRole, err := s.roleRepository.FindByName(normalizedName)
	if err == nil && existingRole.ID != 0 {
		return models.Role{}, models.ErrRoleNameExists
	}

	role := models.Role{
		Name: normalizedName,
	}
	return s.roleRepository.Create(role)
}

// GetRoles retrieves all roles
func (s *roleService) GetRoles() ([]models.Role, error) {
	return s.roleRepository.FindAll()
}

// GetRoleByID retrieves a role by ID
func (s *roleService) GetRoleByID(id uint) (models.Role, error) {
	if id == 0 {
		return models.Role{}, models.ErrRoleIDInvalid
	}
	return s.roleRepository.FindByID(id)
}

// GetRoleByName retrieves a role by name
func (s *roleService) GetRoleByName(name string) (models.Role, error) {
	if name == "" {
		return models.Role{}, models.ErrRoleNameRequired
	}
	return s.roleRepository.FindByName(normalizeRole(name))
}

// UpdateRole updates a role (Super Admin only)
func (s *roleService) UpdateRole(id uint, input dto.CreateRoleInput) (models.Role, error) {
	if id == 0 {
		return models.Role{}, models.ErrRoleIDInvalid
	}

	if input.Name == "" {
		return models.Role{}, models.ErrRoleNameRequired
	}

	role, err := s.roleRepository.FindByID(id)
	if err != nil {
		return models.Role{}, models.ErrRoleNotFound
	}

	normalizedName := normalizeRole(input.Name)

	// Check if new name already exists (excluding current role)
	existingRole, err := s.roleRepository.FindByName(normalizedName)
	if err == nil && existingRole.ID != 0 && existingRole.ID != id {
		return models.Role{}, models.ErrRoleNameExists
	}

	role.Name = normalizedName
	return s.roleRepository.Update(role)
}

// DeleteRole deletes a role (Super Admin only)
func (s *roleService) DeleteRole(id uint) error {
	if id == 0 {
		return models.ErrRoleIDInvalid
	}

	// Verify role exists before deleting
	_, err := s.roleRepository.FindByID(id)
	if err != nil {
		return models.ErrRoleNotFound
	}

	return s.roleRepository.Delete(id)
}
