package dto

type CreateTripInput struct {
	Destination string `json:"destination" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"` // format YYYY-MM-DD
	EndDate     string `json:"end_date" binding:"required"`
}