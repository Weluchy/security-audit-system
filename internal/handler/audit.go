package handler

import (
	"encoding/json"
	"net/http"
	"security-audit-system/internal/model"
	"security-audit-system/internal/service"
)

type AuditHandler struct {
	service *service.AuditService
}

func NewAuditHandler(service *service.AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

func (h *AuditHandler) HandleAudit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var data model.AuditRequest
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Неверный json", http.StatusBadRequest)
			return
		}
		if err := h.service.CreateEvent(data); err != nil {
			http.Error(w, "Ошибка записи на сервере", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)

	case http.MethodGet:
		data, err := h.service.GetEvents()
		if err != nil {
			http.Error(w, "Ошибка записи на сервере", http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Ошибка записи на сервере", http.StatusInternalServerError)
			return
		}
	}
}
