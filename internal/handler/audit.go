package handler

import (
	"net/http"
	"security-audit-system/internal/model"
	"security-audit-system/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuditHandler struct {
	service *service.AuditService
}

func NewAuditHandler(service *service.AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

func (h *AuditHandler) Create(c *gin.Context) {
	var req model.AuditRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateEvent(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, req)

}

func (h *AuditHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}

func (h *AuditHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	strId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data, err := h.service.GetEventByID(strId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, data)

}

func (h *AuditHandler) GetStats(c *gin.Context) {

	data, err := h.service.GetActionCountPerUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}
