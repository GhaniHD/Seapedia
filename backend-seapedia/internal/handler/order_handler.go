package handler

import (
	"net/http"

	"backend-seapedia/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct{ svc service.OrderService }

func NewOrderHandler(svc service.OrderService) *OrderHandler { return &OrderHandler{svc: svc} }

func (h *OrderHandler) ListMine(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	orders, err := h.svc.ListBuyerOrders(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func (h *OrderHandler) GetMyDetail(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id order tidak valid"})
		return
	}
	order, err := h.svc.GetBuyerOrderDetail(c.Request.Context(), userID, orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": order})
}

func (h *OrderHandler) SpendingReport(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	report, err := h.svc.BuyerSpendingReport(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": report})
}

func (h *OrderHandler) ListIncoming(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	orders, err := h.svc.ListSellerOrders(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

func (h *OrderHandler) ProcessOrder(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id order tidak valid"})
		return
	}
	if err := h.svc.ProcessOrder(c.Request.Context(), userID, orderID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order berhasil diproses, sekarang tersedia untuk driver"})
}

func (h *OrderHandler) IncomeReport(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	report, err := h.svc.SellerIncomeReport(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": report})
}
