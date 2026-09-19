package main

import (
	"net/http"
	"security-audit-system/internal/handler"
	"security-audit-system/internal/model"
	"security-audit-system/internal/repository"
	"security-audit-system/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost port=5432 user=myuser password=mypassword dbname=audit_logs sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&model.Event{}); err != nil {
		panic(err)
	}
	repo := repository.NewAuditRepo(db)
	service := service.NewAuditService(repo)
	handler := handler.NewAuditHandler(service)

	http.HandleFunc("/audit", handler.HandleAudit)
	http.ListenAndServe(":8080", nil)
}
