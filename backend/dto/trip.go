package dto

type CreateTripInput struct {
	Destination string `form:"destination" binding:"required,min=3,max=255" example:"Jakarta"`
	StartDate   string `form:"start_date" binding:"required" example:"2024-01-15"` // format YYYY-MM-DD
	EndDate     string `form:"end_date" binding:"required" example:"2024-01-20"`
	Description string `form:"description" binding:"required,min=5" example:"Perjalanan dinas untuk audit"`
	Initiator   string `form:"initiator" binding:"required,min=3,max=100" example:"John Doe"`
	Summary     string `form:"summary" binding:"omitempty,max=500" example:"Ringkasan perjalanan dinas"`
	Nomor_Surat string `form:"nomor_surat" binding:"required,min=2,max=50" example:"DIN/2024/001"`
}

type UpdateTripInput struct {
	Destination string `json:"destination" form:"destination" binding:"omitempty,min=3,max=255" example:"Jakarta"`
	StartDate   string `json:"start_date" form:"start_date" binding:"omitempty" example:"2024-01-15"`
	EndDate     string `json:"end_date" form:"end_date" binding:"omitempty" example:"2024-01-20"`
	Description string `json:"description" form:"description" binding:"omitempty,min=5" example:"Perjalanan dinas untuk audit"`
	Initiator   string `json:"initiator" form:"initiator" binding:"omitempty,min=3,max=100" example:"John Doe"`
	Summary     string `json:"summary" form:"summary" binding:"omitempty,max=500" example:"Ringkasan perjalanan dinas"`
	Nomor_Surat string `json:"nomor_surat" form:"nomor_surat" binding:"omitempty,min=2,max=50" example:"DIN/2024/001"`
}

type UpdateClaimInput struct {
	AttachmentPath string `json:"attachment_path"`
	Notes          string `json:"notes"`
}

type ManagerApprovalInput struct {
	Status string `json:"status" binding:"required"`
	Reason string `json:"reason"`
}

type FinancialReviewInput struct {
	Status string `json:"status" binding:"required"`
	Notes  string `json:"notes"`
}

type FinancialDisbursementInput struct {
	Amount      float64 `json:"amount" binding:"required"`
	ReferenceID string  `json:"reference_id" binding:"required"`
	Notes       string  `json:"notes"`
}

type TripResponse struct {
	ID              uint        `json:"id"`
	UserID          uint        `json:"user_id"`
	User            interface{} `json:"user"`
	Description     string      `json:"description"`
	Destination     string      `json:"destination"`
	StartDate       string      `json:"start_date"`
	EndDate         string      `json:"end_date"`
	Initiator       string      `json:"initiator"`
	Summary         string      `json:"summary,omitempty"`
	Nomor_Surat     string      `json:"nomor_surat"`
	Status          string      `json:"status"`
	AttachmentPath  string      `json:"attachment_path,omitempty"`
	AttachmentPaths string      `json:"attachment_paths,omitempty"`
	Notes           string      `json:"notes,omitempty"`
	CreatedAt       string      `json:"created_at"`
	UpdatedAt       string      `json:"updated_at"`
}
