package handler

import (
	"net/http"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletHandler struct{ svc service.WalletService }

func NewWalletHandler(svc service.WalletService) *WalletHandler { return &WalletHandler{svc: svc} }

func (h *WalletHandler) Topup(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	var req dto.TopupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid: " + err.Error()})
		return
	}
	w, err := h.svc.Topup(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": w})
}

func (h *WalletHandler) GetBalance(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	w, err := h.svc.GetBalance(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": w})
}

func (h *WalletHandler) ListTransactions(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	txs, err := h.svc.ListTransactions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": txs})
}

func (h *WalletHandler) AddAddress(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	var req dto.UpsertAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid: " + err.Error()})
		return
	}
	a, err := h.svc.AddAddress(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": a})
}

func (h *WalletHandler) ListAddresses(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	addrs, err := h.svc.ListAddresses(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": addrs})
}
