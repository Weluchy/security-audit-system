package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5432 user=myuser password=mypassword dbname=audit_logs sslmode=disable")
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

	http.HandleFunc("/audit", handleAudit)
	http.ListenAndServe(":8080", nil)
}

type AuditRequest struct {
	UserID int    `json:"user_id"`
	Action string `json:"action"`
}

func handleAudit(w http.ResponseWriter, r *http.Request) {

	var auditReq AuditRequest

	switch r.Method {
	case http.MethodPost:
		if err := json.NewDecoder(r.Body).Decode(&auditReq); err != nil {
			http.Error(w, "Неверный json", http.StatusBadRequest)
			return
		}
		if _, err := db.Exec(`insert into events (user_id, action) values ($1, $2)`, auditReq.UserID, auditReq.Action); err != nil {
			http.Error(w, "Не получилось добавить данные", http.StatusInternalServerError)
			return
		}
	case http.MethodGet:

		rows, err := db.Query(`select user_id, action from events`)
		if err != nil {
			http.Error(w, "Не получилось просмотреть данные", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var slice []AuditRequest
		for rows.Next() {
			var audit AuditRequest
			if err := rows.Scan(&audit.UserID, &audit.Action); err != nil {
				http.Error(w, "Ошибка rows", http.StatusInternalServerError)
				return
			}
			slice = append(slice, audit)
		}

		if err = rows.Err(); err != nil {
			http.Error(w, "Ошибка итерации", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(slice); err != nil {
			http.Error(w, "Ошибка кодировки json", http.StatusInternalServerError)
			return
		}

	}

}
