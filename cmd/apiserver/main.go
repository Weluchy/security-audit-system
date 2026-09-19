package main

import (
	"security-audit-system/internal/handler"
	"security-audit-system/internal/model"
	"security-audit-system/internal/repository"
	"security-audit-system/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost port=5433 user=myuser password=mypassword dbname=audit_logs sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&model.Event{}); err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "",
		DB:       0,
	})
	if err := rdb.Ping().Err(); err != nil {
		panic(err)
	}

	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "audit_events",
		Balancer: &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()

	repo := repository.NewAuditRepo(db, rdb)
	service := service.NewAuditService(repo, kafkaWriter)
	handler := handler.NewAuditHandler(service)

	router := gin.Default()
	router.GET("/audit", handler.GetAll)
	router.GET("/audit/:id", handler.GetByID)
	router.POST("/audit", handler.Create)
	router.GET("/audit/count", handler.GetStats)
	router.Run(":8080")
}
