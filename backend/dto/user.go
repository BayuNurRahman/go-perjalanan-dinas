package dto

type UserResponse struct {
	ID           uint   `json:"id" example:"10"`
	Name         string `json:"name" example:"Budi Santoso"`
	Email        string `json:"email" example:"budi@example.com"`
	Role         string `json:"role" example:"EMPLOYEE"`
	DepartmentID *uint  `json:"department_id,omitempty" example:"1"`
}

type UpdateUserInput struct {
	Name         string `json:"name,omitempty" example:"Budi Santoso"`
	Email        string `json:"email,omitempty" binding:"email" example:"budi@example.com"`
	RoleID       *uint  `json:"role_id,omitempty" example:"2"`
	DepartmentID *uint  `json:"department_id,omitempty" example:"1"`
}
