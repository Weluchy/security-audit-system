package main

import (
	"context"
	"log"
	"security-audit-system/internal/repository"
	"security-audit-system/internal/worker"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("Запуск Воркера Kafka Consumer...")

	dsn := "host=localhost port=5433 user=myuser password=mypassword dbname=audit_logs sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}
	repo := repository.NewAuditRepo(db, nil)
	consumer := worker.NewAuditConsumer(repo)
	consumer.Start(context.Background())
}
