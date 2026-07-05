package handler

import (
	"net/http"

	"backend-seapedia/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeliveryHandler struct{ svc service.DeliveryService }

func NewDeliveryHandler(svc service.DeliveryService) *DeliveryHandler { return &DeliveryHandler{svc: svc} }

func (h *DeliveryHandler) FindJobs(c *gin.Context) {
	jobs, err := h.svc.FindAvailableJobs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": jobs})
}

func (h *DeliveryHandler) GetDetail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id job tidak valid"})
		return
	}
	job, err := h.svc.GetJobDetail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": job})
}

func (h *DeliveryHandler) TakeJob(c *gin.Context) {
	driverID, _ := uuid.Parse(c.GetString("user_id"))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id job tidak valid"})
		return
	}
	if err := h.svc.TakeJob(c.Request.Context(), driverID, id); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "job berhasil diambil"})
}

func (h *DeliveryHandler) CompleteJob(c *gin.Context) {
	driverID, _ := uuid.Parse(c.GetString("user_id"))
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id job tidak valid"})
		return
	}
	if err := h.svc.CompleteJob(c.Request.Context(), driverID, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "job dikonfirmasi selesai"})
}

func (h *DeliveryHandler) MyJobs(c *gin.Context) {
	driverID, _ := uuid.Parse(c.GetString("user_id"))
	jobs, err := h.svc.MyJobs(c.Request.Context(), driverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": jobs})
}

func (h *DeliveryHandler) MyEarnings(c *gin.Context) {
	driverID, _ := uuid.Parse(c.GetString("user_id"))
	earnings, err := h.svc.MyEarnings(c.Request.Context(), driverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": earnings})
}
