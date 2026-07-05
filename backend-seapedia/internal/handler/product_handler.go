package handler

import (
	"net/http"

	"backend-seapedia/internal/dto"
	"backend-seapedia/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct{ svc service.ProductService }

func NewProductHandler(svc service.ProductService) *ProductHandler { return &ProductHandler{svc: svc} }

func (h *ProductHandler) Create(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	var req dto.UpsertProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid: " + err.Error()})
		return
	}
	p, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": p})
}

func (h *ProductHandler) Update(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id produk tidak valid"})
		return
	}
	var req dto.UpsertProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data tidak valid: " + err.Error()})
		return
	}
	p, err := h.svc.Update(c.Request.Context(), userID, productID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func (h *ProductHandler) Delete(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id produk tidak valid"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), userID, productID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "produk berhasil dihapus"})
}

func (h *ProductHandler) ListMine(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	products, err := h.svc.ListMine(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": products})
}

func (h *ProductHandler) ListPublic(c *gin.Context) {
	products, err := h.svc.ListPublic(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": products})
}

func (h *ProductHandler) GetPublicDetail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id produk tidak valid"})
		return
	}
	p, err := h.svc.GetPublicDetail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}
