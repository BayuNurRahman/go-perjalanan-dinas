package repository

import (
	"gorm.io/gorm"
	"go-perjalanan-dinas/models"
)

// 1. Define the UserRepository interface with methods for creating a user and finding a user by email.
type UserRepository interface {
	Create(user models.User) (models.User, error)
	FindByEmail(email string) (models.User, error)
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
	return user, err
}

// 5. Implement the FindByEmail method to retrieve a user record from the database based on the provided email address and return the user along with any error encountered.
func (r *userRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}