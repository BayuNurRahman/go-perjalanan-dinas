package models

type Department struct {
	ID   uint   `gorm:"primaryKey; autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(100);not null" json:"name"` // Information Technology, Human Resources, Finance, Marketing, etc.
	Code string `gorm:"type:varchar(100);not null" json:"code"` // IT, HR, FIN, MKT, etc.
}
