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

func (rep *AuditRepo) GetByID(id int) (model.AuditRequest, error) {
	var event model.Event
	if err := rep.db.First(&event, id).Error; err != nil {
		return model.AuditRequest{}, err
	}
	return model.AuditRequest{
		UserID: event.UserID,
		Action: event.Action,
	}, nil
}

func (rep *AuditRepo) GetActionCountPerUser() (map[int]int, error) {
	type Result struct {
		UserID int
		Count  int
	}
	var resultsQuery []Result

	if err := rep.db.Raw("select user_id, count(*) as count from events group by user_id").Scan(&resultsQuery).Error; err != nil {
		return nil, err
	}

	results := make(map[int]int)
	for _, v := range resultsQuery {
		results[v.UserID] = v.Count
	}
	return results, nil
}
