package service

import (
	"security-audit-system/internal/model"
	"security-audit-system/internal/repository"
)

type AuditService struct {
	repo *repository.AuditRepo
}

func NewAuditService(repo *repository.AuditRepo) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) CreateEvent(req model.AuditRequest) error {
	return s.repo.Save(req)
}

func (s *AuditService) GetEvents() ([]model.AuditRequest, error) {
	return s.repo.GetAll()
}
