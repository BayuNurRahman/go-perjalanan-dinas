package dto

// WebResponse represents the standard API response envelope used across the project.
type WebResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty" example:"Operasi berhasil"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse remains as a compatibility alias for Swagger/docs.
type SuccessResponse = WebResponse

// ErrorResponse remains as a compatibility alias for Swagger/docs.
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Terjadi kesalahan"`
}
