package dto

type SubmitClaimInput struct {
	TripID          uint    `json:"trip_id" form:"trip_id" binding:"required"`
	Title           string  `json:"title" form:"title" binding:"required,min=3,max=100"`
	Description     string  `json:"description" form:"description" binding:"required,min=3"`
	Amount          float64 `json:"amount" form:"amount" binding:"required,gt=0"`
	TransactionDate string  `json:"transaction_date" form:"transaction_date" binding:"required"`
}

type ReviewClaimInput struct {
	Status         string `json:"status" form:"status" binding:"required"`
	RejectedReason string `json:"rejected_reason" form:"rejected_reason"`
}
