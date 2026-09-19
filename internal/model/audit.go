package model

type AuditRequest struct {
	UserID int    `json:"user_id"`
	Action string `json:"action"`
}

type Event struct {
	ID     uint `gorm:"primaryKey"`
	UserID int
	Action string
}
