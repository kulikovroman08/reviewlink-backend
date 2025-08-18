package controller

import (
	"net/http"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
)

func (h *Application) CreatePlace(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var req dto.CreatePlaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	place := model.Place{
		Name:    req.Name,
		Address: req.Address,
	}

	resp, err := h.Service.CreatePlace(c.Request.Context(), place)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create place"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
