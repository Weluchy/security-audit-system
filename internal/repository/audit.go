package repository

import (
	"database/sql"
	"security-audit-system/internal/model"
)

type AuditRepo struct {
	db *sql.DB
}

func NewAuditRepo(db *sql.DB) *AuditRepo {
	return &AuditRepo{db: db}
}

func (rep *AuditRepo) Save(req model.AuditRequest) error {
	if _, err := rep.db.Exec(`insert into events (user_id, action) values ($1, $2)`, req.Action, req.Action); err != nil {
		return err
	}
	return nil
}

func (req *AuditRepo) GetAll() ([]model.AuditRequest, error) {
	var resp model.AuditRequest
	var result []model.AuditRequest
	rows, err := req.db.Query(`select user_id, action from events`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&resp.UserID, &resp.Action); err != nil {
			return nil, err
		}
		result = append(result, resp)
	}
	return result, nil
}
