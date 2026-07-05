package handler

import (
	"net/http"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CheckoutHandler struct{ svc service.CheckoutService }

func NewCheckoutHandler(svc service.CheckoutService) *CheckoutHandler { return &CheckoutHandler{svc: svc} }

func (h *CheckoutHandler) Checkout(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	var req dto.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid: " + err.Error()})
		return
	}
	summary, err := h.svc.Checkout(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": summary})
}
