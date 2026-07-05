package handler

import (
	"net/http"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartHandler struct{ svc service.CartService }

func NewCartHandler(svc service.CartService) *CartHandler { return &CartHandler{svc: svc} }

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	cart, err := h.svc.GetCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cart})
}

func (h *CartHandler) AddItem(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	var req dto.AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid: " + err.Error()})
		return
	}
	cart, err := h.svc.AddItem(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cart})
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id tidak valid"})
		return
	}
	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid: " + err.Error()})
		return
	}
	cart, err := h.svc.UpdateItem(c.Request.Context(), userID, productID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cart})
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id tidak valid"})
		return
	}
	cart, err := h.svc.RemoveItem(c.Request.Context(), userID, productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cart})
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	if err := h.svc.ClearCart(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cart berhasil dikosongkan"})
}
