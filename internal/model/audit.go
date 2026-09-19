package model

type AuditRequest struct {
	UserID int    `json:"user_id"`
	Action string `json:"action"`
}
