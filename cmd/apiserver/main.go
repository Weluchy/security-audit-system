package main

import (
	"security-audit-system/internal/handler"
	"security-audit-system/internal/model"
	"security-audit-system/internal/repository"
	"security-audit-system/internal/service"

	"github.com/gin-gonic/gin"
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

	router := gin.Default()
	router.GET("/audit", handler.GetAll)
	router.GET("/audit/:id", handler.GetByID)
	router.POST("/audit", handler.Create)
	router.GET("/audit/count", handler.GetStats)
	router.Run(":8080")
}
