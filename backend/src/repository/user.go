package repository

import (
	"go-perjalanan-dinas/models"

	"gorm.io/gorm"
)

// 1. Define the UserRepository interface with methods for creating a user and finding a user by email.
type UserRepository interface {
	Create(user models.User) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindAll() ([]models.User, error)
	FindByID(id uint) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(id uint) error
}

// 2. Implement the UserRepository interface with a struct that holds a reference to the GORM database connection.
type userRepository struct {
	db *gorm.DB
}

// 3. Create a constructor function that initializes the userRepository with a GORM database connection and returns it as a UserRepository interface.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// 4. Implement the Create method to insert a new user record into the database and return the created user along with any error encountered.
func (r *userRepository) Create(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	// Query back the created user with role relationship to ensure role name is populated
	// Don't use association on just-created user; query from database instead
	return r.FindByID(user.ID)
}

// 5. Implement the FindByEmail method to retrieve a user record from the database based on the provided email address and return the user along with any error encountered.
func (r *userRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Preload("Department").Where("email = ?", email).First(&user).Error
	if err != nil {
		return user, err
	}
	return r.loadRoleName(user)
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Department").Preload("RoleData").Find(&users).Error
	if err != nil {
		return users, err
	}
	// Populate role names for all users
	for i := range users {
		users[i] = r.populateRoleName(users[i])
	}
	return users, nil
}

func (r *userRepository) FindByID(id uint) (models.User, error) {
	var user models.User
	err := r.db.Preload("Department").Preload("RoleData").First(&user, id).Error
	if err != nil {
		return user, err
	}
	return r.populateRoleName(user), nil
}

func (r *userRepository) Update(user models.User) (models.User, error) {
	err := r.db.Save(&user).Error
	if err != nil {
		return user, err
	}
	return r.FindByID(user.ID)
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// loadRoleName loads role relationship and populates denormalized role name
func (r *userRepository) loadRoleName(user models.User) (models.User, error) {
	if user.RoleID != nil {
		err := r.db.Model(&user).Association("RoleData").Find(&user.RoleData)
		if err != nil {
			return user, nil
		}
	}
	return r.populateRoleName(user), nil
}

// populateRoleName copies role name from RoleData to denormalized Role field
func (r *userRepository) populateRoleName(user models.User) models.User {
	if user.RoleData.Name != "" {
		user.Role = user.RoleData.Name
	}
	return user
}
