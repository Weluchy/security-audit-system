package main

import (
	"database/sql"
	"net/http"
	"security-audit-system/internal/handler"
	"security-audit-system/internal/repository"
	"security-audit-system/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=myuser password=mypassword dbname=audit_logs sslmode=disable")
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	defer db.Close()

	if _, err := db.Exec(`create table if not exists events (
	id generated always as identity primary key,
	user_id int,
	action varchar(255))`); err != nil {
		panic(err)
	}
	repo := repository.NewAuditRepo(db)
	service := service.NewAuditService(repo)
	handler := handler.NewAuditHandler(service)

	http.HandleFunc("/audit", handler.HandleAudit)
	http.ListenAndServe(":8080", nil)
}
