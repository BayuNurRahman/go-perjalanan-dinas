package dto

type RegisterInput struct {
	Name         string `json:"name" binding:"required" example:"Budi Santoso"`
	Email        string `json:"email" binding:"required,email" example:"budi@example.com"`
	Password     string `json:"password" binding:"required,min=6" example:"password123"`
	RoleID       uint   `json:"role_id" binding:"required" example:"2"`
	DepartmentID uint   `json:"department_id" binding:"required" example:"1"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email" example:"budi@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}
