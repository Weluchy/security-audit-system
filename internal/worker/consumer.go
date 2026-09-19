package worker

import (
	"context"
	"encoding/json"
	"log"
	"security-audit-system/internal/model"
	"security-audit-system/internal/repository"

	"github.com/segmentio/kafka-go"
)

type AuditConsumer struct {
	reader *kafka.Reader
	repo   *repository.AuditRepo
}

func NewAuditConsumer(repo *repository.AuditRepo) *AuditConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "audit_events",
		GroupID: "audit_writer_group",
	})
	return &AuditConsumer{repo: repo, reader: r}
}

func (c *AuditConsumer) Start(ctx context.Context) {
	log.Println("Consumer запущен...")
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Println("Ошибка ", err)
			continue
		}
		var req model.AuditRequest
		if err := json.Unmarshal(msg.Value, &req); err != nil {
			log.Println("Ошибка парсинга JSON:", err)
			continue
		}

		if err := c.repo.Save(req); err != nil {
			log.Println("Ошибка ", err)
			continue
		}
		log.Println("Событие успешно сохранено в БД!")
	}
}
