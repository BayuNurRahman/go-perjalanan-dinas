package service

import (
	"errors"
	"testing"

	"go-perjalanan-dinas/dto"
	"go-perjalanan-dinas/models"

	"gorm.io/gorm"
)

type fakeDepartmentRepo struct {
	departments []models.Department
}

func (f *fakeDepartmentRepo) Create(department models.Department) (models.Department, error) {
	f.departments = append(f.departments, department)
	return department, nil
}

func (f *fakeDepartmentRepo) FindAll() ([]models.Department, error) {
	return f.departments, nil
}

func (f *fakeDepartmentRepo) FindByID(id uint) (models.Department, error) {
	for _, dept := range f.departments {
		if dept.ID == id {
			return dept, nil
		}
	}
	return models.Department{}, nil
}

func (f *fakeDepartmentRepo) Update(department models.Department) (models.Department, error) {
	for i, dept := range f.departments {
		if dept.ID == department.ID {
			f.departments[i] = department
			return department, nil
		}
	}
	f.departments = append(f.departments, department)
	return department, nil
}

func (f *fakeDepartmentRepo) Delete(id uint) error {
	filtered := f.departments[:0]
	for _, dept := range f.departments {
		if dept.ID != id {
			filtered = append(filtered, dept)
		}
	}
	f.departments = filtered
	return nil
}

type fakeUserRepo struct {
	users []models.User
}

func (f *fakeUserRepo) Create(user models.User) (models.User, error) {
	f.users = append(f.users, user)
	return user, nil
}

func (f *fakeUserRepo) FindByEmail(email string) (models.User, error) {
	for _, u := range f.users {
		if u.Email == email {
			return u, nil
		}
	}
	return models.User{}, errors.New("not found")
}

func (f *fakeUserRepo) FindAll() ([]models.User, error) {
	return f.users, nil
}

func (f *fakeUserRepo) FindByID(id uint) (models.User, error) {
	for _, u := range f.users {
		if u.ID == id {
			return u, nil
		}
	}
	return models.User{}, nil
}

func (f *fakeUserRepo) Update(user models.User) (models.User, error) {
	for i, existing := range f.users {
		if existing.ID == user.ID {
			f.users[i] = user
			return user, nil
		}
	}
	f.users = append(f.users, user)
	return user, nil
}

func (f *fakeUserRepo) Delete(id uint) error {
	filtered := f.users[:0]
	for _, user := range f.users {
		if user.ID != id {
			filtered = append(filtered, user)
		}
	}
	f.users = filtered
	return nil
}

func TestDepartmentServiceCreateDepartment(t *testing.T) {
	repo := &fakeDepartmentRepo{}
	svc := NewDepartmentService(repo)

	dept, err := svc.CreateDepartment("SUPER_ADMIN", dto.CreateDepartmentInput{Name: "IT", Code: "it"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if dept.Name != "IT" || dept.Code != "it" {
		t.Fatalf("unexpected department payload: %+v", dept)
	}
}

func TestUserServiceUpdateUser(t *testing.T) {
	managerRoleID := uint(1)
	repo := &fakeUserRepo{users: []models.User{{Model: gorm.Model{ID: 1}, Name: "Ana", Email: "ana@example.com", RoleID: &managerRoleID, Role: "MANAGER"}}}
	svc := NewUserService(repo)

	newRoleID := uint(2)
	updated, err := svc.UpdateUser(1, dto.UpdateUserInput{Name: "Ana Updated", RoleID: &newRoleID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Name != "Ana Updated" || updated.RoleID == nil || *updated.RoleID != newRoleID {
		t.Fatalf("unexpected updated user: %+v", updated)
	}
}

func TestRegisterAcceptsSuperAdminRole(t *testing.T) {
	roleID := uint(1)
	repo := &fakeUserRepo{users: []models.User{}}
	svc := &authService{userRepository: repo}

	user, err := svc.Register(dto.RegisterInput{Name: "Admin", Email: "super@example.com", Password: "password123", RoleID: roleID, DepartmentID: 1})
	if err != nil {
		t.Fatalf("expected super admin role to be accepted, got %v", err)
	}
	if user.RoleID == nil || *user.RoleID != roleID {
		t.Fatalf("expected RoleID %d, got %+v", roleID, user.RoleID)
	}
}
