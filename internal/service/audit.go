package service

import (
	"context"
	"encoding/json"
	"fmt"
	"security-audit-system/internal/model"

	"github.com/segmentio/kafka-go"
)

type AuditRepository interface {
	GetByID(id int) (model.AuditRequest, error)
	GetAll() ([]model.AuditRequest, error)
	GetActionCountPerUser() (map[int]int, error)
}
type AuditService struct {
	repo        AuditRepository
	kafkaWriter *kafka.Writer
}

func NewAuditService(repo AuditRepository, kafkaWriter *kafka.Writer) *AuditService {
	return &AuditService{repo: repo, kafkaWriter: kafkaWriter}
}

func (s *AuditService) CreateEvent(req model.AuditRequest) error {

	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("user-%d", req.UserID)),
		Value: bytes,
	}

	err = s.kafkaWriter.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}

	return nil

}

func (s *AuditService) GetEvents() ([]model.AuditRequest, error) {
	return s.repo.GetAll()
}
func (s *AuditService) GetEventByID(id int) (model.AuditRequest, error) {
	return s.repo.GetByID(id)
}
func (s *AuditService) GetActionCountPerUser() (map[int]int, error) {
	return s.repo.GetActionCountPerUser()
}
