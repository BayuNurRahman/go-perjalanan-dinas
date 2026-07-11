package service

import (
	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"
	"go-perjalanan-dinas/src/repository"
)

type DepartmentService interface {
	CreateDepartment(userRole string, input dto.CreateDepartmentInput) (models.Department, error)
	GetDepartments() ([]models.Department, error)
	GetDepartmentByID(id uint) (models.Department, error)
	UpdateDepartment(id uint, input dto.CreateDepartmentInput) (models.Department, error)
	DeleteDepartment(id uint) error
}

type departmentService struct {
	departmentRepository repository.DepartmentRepository
}

func NewDepartmentService(repo repository.DepartmentRepository) DepartmentService {
	return &departmentService{repo}
}

// CreateDepartment creates a new department (Super Admin only)
func (s *departmentService) CreateDepartment(userRole string, input dto.CreateDepartmentInput) (models.Department, error) {
	if normalizeRole(userRole) != "SUPER_ADMIN" {
		return models.Department{}, models.ErrSuperAdminOnly
	}

	if input.Name == "" || input.Code == "" {
		return models.Department{}, models.ErrDepartmentInputRequired
	}

	department := models.Department{
		Name: input.Name,
		Code: input.Code,
	}
	return s.departmentRepository.Create(department)
}

// GetDepartments retrieves all departments
func (s *departmentService) GetDepartments() ([]models.Department, error) {
	return s.departmentRepository.FindAll()
}

// GetDepartmentByID retrieves a department by ID
func (s *departmentService) GetDepartmentByID(id uint) (models.Department, error) {
	if id == 0 {
		return models.Department{}, models.ErrDepartmentIDInvalid
	}
	return s.departmentRepository.FindByID(id)
}

// UpdateDepartment updates a department (Super Admin only)
func (s *departmentService) UpdateDepartment(id uint, input dto.CreateDepartmentInput) (models.Department, error) {
	if id == 0 {
		return models.Department{}, models.ErrDepartmentIDInvalid
	}

	if input.Name == "" && input.Code == "" {
		return models.Department{}, models.ErrDepartmentUpdateRequired
	}

	department, err := s.departmentRepository.FindByID(id)
	if err != nil {
		return models.Department{}, models.ErrDepartmentNotFound
	}

	if input.Name != "" {
		department.Name = input.Name
	}
	if input.Code != "" {
		department.Code = input.Code
	}

	return s.departmentRepository.Update(department)
}

// DeleteDepartment deletes a department (Super Admin only)
func (s *departmentService) DeleteDepartment(id uint) error {
	if id == 0 {
		return models.ErrDepartmentIDInvalid
	}

	// Verify department exists before deleting
	_, err := s.departmentRepository.FindByID(id)
	if err != nil {
		return models.ErrDepartmentNotFound
	}

	return s.departmentRepository.Delete(id)
}
