package dto

// CreateDepartmentInput adalah struktur untuk input pembuatan departemen baru
type CreateDepartmentInput struct {
	Name string `json:"name" binding:"required,min=3,max=100" example:"Information Technology"` // Nama departemen
	Code string `json:"code" binding:"required,min=2,max=10" example:"IT"`                      // Kode departemen (singkat)
}

// UpdateDepartmentInput adalah struktur untuk input update departemen
type UpdateDepartmentInput struct {
	Name string `json:"name" binding:"omitempty,min=3,max=100" example:"Information Technology"` // Nama departemen (opsional)
	Code string `json:"code" binding:"omitempty,min=2,max=10" example:"IT"`                      // Kode departemen (opsional)
}

// DepartmentResponse adalah struktur untuk response departemen di API
type DepartmentResponse struct {
	ID   uint   `json:"id" example:"1"`                        // ID departemen
	Name string `json:"name" example:"Information Technology"` // Nama departemen
	Code string `json:"code" example:"IT"`                     // Kode departemen
}
