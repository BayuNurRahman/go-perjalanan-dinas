package dto

type CreateRoleInput struct {
	Name string `json:"name" binding:"required" example:"MANAGER"`
}

type RoleResponse struct {
	ID   uint   `json:"id" example:"3"`
	Name string `json:"name" example:"MANAGER"`
}
