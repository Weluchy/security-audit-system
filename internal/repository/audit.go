package repository

import (
	"security-audit-system/internal/model"

	"gorm.io/gorm"
)

type AuditRepo struct {
	db *gorm.DB
}

func NewAuditRepo(db *gorm.DB) *AuditRepo {
	return &AuditRepo{db: db}
}

func (rep *AuditRepo) Save(req model.AuditRequest) error {
	event := model.Event{
		UserID: req.UserID,
		Action: req.Action,
	}
	if err := rep.db.Create(&event).Error; err != nil {
		return err
	}
	return nil
}

func (rep *AuditRepo) GetAll() ([]model.AuditRequest, error) {
	var events []model.Event
	var result []model.AuditRequest
	if err := rep.db.Find(&events).Error; err != nil {
		return nil, err
	}
	for _, event := range events {
		result = append(result, model.AuditRequest{
			UserID: event.UserID,
			Action: event.Action,
		})
	}
	return result, nil
}
