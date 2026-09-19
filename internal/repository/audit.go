package repository

import (
	"encoding/json"
	"fmt"
	"security-audit-system/internal/model"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type AuditRepo struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewAuditRepo(db *gorm.DB, redis *redis.Client) *AuditRepo {
	return &AuditRepo{db: db, redis: redis}
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

	key := fmt.Sprintf("event:%v", id)
	val, err := rep.redis.Get(key).Result()
	if err == nil {
		err = json.Unmarshal([]byte(val), &event)

		if err != nil {
			return model.AuditRequest{}, err
		}

		return model.AuditRequest{
			UserID: event.UserID,
			Action: event.Action,
		}, nil

	} else if err == redis.Nil {
		if err := rep.db.First(&event, id).Error; err != nil {
			return model.AuditRequest{}, err
		}
		bytes, err := json.Marshal(&event)
		if err != nil {
			return model.AuditRequest{}, err
		}

		rep.redis.Set(key, bytes, 5*time.Minute)

		return model.AuditRequest{
			UserID: event.UserID,
			Action: event.Action,
		}, nil
	}

	return model.AuditRequest{}, err
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
