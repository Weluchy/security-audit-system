package service

import (
	"errors"
	"security-audit-system/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetByID(id int) (model.AuditRequest, error) {
	args := m.Called(id)
	return args.Get(0).(model.AuditRequest), args.Error(1)
}
func (m *MockRepository) GetAll() ([]model.AuditRequest, error)       { return nil, nil }
func (m *MockRepository) GetActionCountPerUser() (map[int]int, error) { return nil, nil }

func TestGetEventByID_Success(t *testing.T) {
	mockRepo := new(MockRepository)

	expectedEvent := model.AuditRequest{UserID: 777, Action: "login"}
	mockRepo.On("GetByID", 5).Return(expectedEvent, nil)

	svc := NewAuditService(mockRepo, nil)

	result, err := svc.GetEventByID(5)

	assert.NoError(t, err)
	assert.Equal(t, 777, result.UserID)
	assert.Equal(t, "login", result.Action)
	mockRepo.AssertExpectations(t)
}

func TestGetEventByID_DatabaseError(t *testing.T) {
	mockRepo := new(MockRepository)

	mockRepo.On("GetByID", 10).Return(model.AuditRequest{}, errors.New("db is down"))

	svc := NewAuditService(mockRepo, nil)
	result, err := svc.GetEventByID(10)

	assert.Error(t, err)
	assert.Equal(t, "db is down", err.Error())
	assert.Equal(t, 0, result.UserID)

	mockRepo.AssertExpectations(t)
}
