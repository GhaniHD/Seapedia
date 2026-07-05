package handler

import (
	"net/http"

	"backend-seapedia/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminSvc   service.AdminService
	overdueSvc service.OverdueService
}

func NewAdminHandler(adminSvc service.AdminService, overdueSvc service.OverdueService) *AdminHandler {
	return &AdminHandler{adminSvc: adminSvc, overdueSvc: overdueSvc}
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	dash, err := h.adminSvc.Dashboard(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": dash})
}

// SimulateNextDay: trigger manual Admin untuk mensimulasikan pergantian hari (Level 6).
// Body opsional: {"days": 1}
func (h *AdminHandler) SimulateNextDay(c *gin.Context) {
	days := 1
	var body struct {
		Days int `json:"days"`
	}
	if err := c.ShouldBindJSON(&body); err == nil && body.Days > 0 {
		days = body.Days
	}
	result, err := h.overdueSvc.SimulateNextDay(c.Request.Context(), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
