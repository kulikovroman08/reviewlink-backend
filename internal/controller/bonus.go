package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/response"
	srvErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"
)

// RedeemBonus godoc
// @Summary Обмен баллов на бонус
// @Description Пользователь может обменять баллы на одно из доступных вознаграждений
// @Security BearerAuth
// @Tags bonuses
// @Accept json
// @Produce json
// @Param request body dto.BonusRedeemRequest true "Данные для получения бонуса (reward_type: free_coffee, free_meal, discount_10)"
// @Success 201 {object} dto.BonusRedeemResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /bonuses/redeem [post]
func (h *Application) RedeemBonus(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var req dto.BonusRedeemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	bonus, err := h.BonusService.RedeemBonus(ctx, userID, req.PlaceID, req.RewardType)
	if err != nil {
		switch {
		case errors.Is(err, srvErrors.ErrNotEnoughPoints):
			ctx.JSON(http.StatusConflict, dto.ErrorResponse{Error: response.ErrNotEnoughPoints})
		case errors.Is(err, srvErrors.ErrInvalidPlaceID):
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidPlaceData})
		case errors.Is(err, srvErrors.ErrBonusCreateFail):
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrInternalError})
		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrInternalError})
		}
		return
	}

	resp := dto.BonusRedeemResponse{
		ID:             bonus.ID.String(),
		PlaceID:        bonus.PlaceID.String(),
		RewardType:     bonus.RewardType,
		RequiredPoints: bonus.RequiredPoints,
		QRToken:        bonus.QRToken,
		IsUsed:         bonus.IsUsed,
		UsedAt:         bonus.UsedAt,
	}
	ctx.JSON(http.StatusCreated, resp)
}

// GetUserBonuses godoc
// @Summary Получить список бонусов пользователя
// @Security BearerAuth
// @Tags bonuses
// @Produce json
// @Success 200 {array}  dto.BonusResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /bonuses [get]
func (h *Application) GetUserBonuses(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	bonuses, err := h.BonusService.GetUserBonuses(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrInternalError})
		return
	}

	var resp []dto.BonusResponse
	for _, b := range bonuses {
		resp = append(resp, dto.BonusResponse{
			ID:             b.ID.String(),
			PlaceID:        b.PlaceID.String(),
			RewardType:     b.RewardType,
			RequiredPoints: b.RequiredPoints,
			QRToken:        b.QRToken,
			IsUsed:         b.IsUsed,
			UsedAt:         b.UsedAt,
		})
	}

	ctx.JSON(http.StatusOK, resp)
}
