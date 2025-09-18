package controller

import (
	"errors"
	"net/http"

	serviceErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"

	"github.com/gin-gonic/gin"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
)

func (a *Application) GenerateTokens(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can generate tokens"})
		return
	}

	var req dto.GenerateTokensRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	resp, err := a.TokenService.GenerateTokens(c.Request.Context(), req.PlaceID, req.Count)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrInvalidPlaceID):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid place id"})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}
