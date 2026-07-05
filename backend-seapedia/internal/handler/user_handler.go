package handler

import (
	"net/http"

	"backend-seapedia/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct{ userService service.UserService }

func NewUserHandler(userService service.UserService) *UserHandler { return &UserHandler{userService: userService} }

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id tidak valid"})
		return
	}
	activeRole := c.GetString("active_role")
	profile, err := h.userService.GetProfile(c.Request.Context(), userID, activeRole)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": profile})
}
