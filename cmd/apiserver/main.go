package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/audit", handleAudit)
	http.ListenAndServe(":8080", nil)
}

type AuditRequest struct {
	UserID int    `json:"user_id"`
	Action string `json:"action"`
}

func handleAudit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Не тот метод", http.StatusMethodNotAllowed)
		return
	}
	var auditReq AuditRequest
	if err := json.NewDecoder(r.Body).Decode(&auditReq); err != nil {
		http.Error(w, "Неверный json", http.StatusBadRequest)
		return
	}

	fmt.Printf("ID: %v, action: %v", auditReq.UserID, auditReq.Action)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{"status":"Декордер отработал"}`))
}
