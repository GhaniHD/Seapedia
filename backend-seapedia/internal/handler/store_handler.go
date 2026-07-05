package handler

import (
	"net/http"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StoreHandler struct{ svc service.StoreService }

func NewStoreHandler(svc service.StoreService) *StoreHandler { return &StoreHandler{svc: svc} }

func (h *StoreHandler) UpsertMyStore(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	var req dto.UpsertStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid: " + err.Error()})
		return
	}
	store, err := h.svc.UpsertMyStore(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": store})
}

func (h *StoreHandler) GetMyStore(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	store, err := h.svc.GetMyStore(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": store})
}

func (h *StoreHandler) ListPublic(c *gin.Context) {
	stores, err := h.svc.ListStores(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stores})
}

func (h *StoreHandler) GetPublicDetail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id toko tidak valid"})
		return
	}
	store, err := h.svc.GetPublicStore(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": store})
}
