package controller

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/response"
	serviceErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"
)

// CreateOwnerRequest godoc
// @Summary      Заявка на владельца заведения
// @Description  Пользователь отправляет заявку стать владельцем заведения из публичного каталога (OSM)
// @Tags         owner_requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CreateOwnerRequest  true  "source/source_id(or sourceId)/name"
// @Success      201   {object}  dto.CreateOwnerRequestResponse
// @Failure      400   {object}  dto.ErrorResponse "invalid input"
// @Failure      401   {object}  dto.ErrorResponse "authentication required"
// @Failure      409   {object}  dto.ErrorResponse "owner request already exists"
// @Failure      500   {object}  dto.ErrorResponse "failed to create owner request"
// @Router       /owner_requests [post]
func (h *Application) CreateOwnerRequest(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: response.ErrUnauthorized})
		return
	}
	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: response.ErrUnauthorized})
		return
	}

	var req dto.CreateOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	source := strings.TrimSpace(req.Source)
	sourceID := strings.TrimSpace(req.SourceID)
	if sourceID == "" {
		sourceID = strings.TrimSpace(req.SourceId)
	}
	name := strings.TrimSpace(req.Name)

	if source == "" || sourceID == "" || name == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	id, err := h.OwnerRequestService.Create(c.Request.Context(), userID, source, sourceID, name)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrOwnerRequestAlreadyExist):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: response.ErrOwnerRequestAlreadyExists})
		case errors.Is(err, serviceErrors.ErrInvalidOwnerRequestData):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrFailedCreateOwnerRequest})
		}
		return
	}

	c.JSON(http.StatusCreated, dto.CreateOwnerRequestResponse{ID: id})
}

// ApproveOwnerRequest godoc
// @Summary      Подтвердить владельца заведения (platform admin)
// @Description  Подтверждает заявку, создаёт/находит place по (source, source_id), назначает owner_id
// @Tags         platform_admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path     string  true  "OwnerRequest ID (uuid)"
// @Param        body  body     dto.ApproveOwnerRequestRequest  false  "comment"
// @Success      200 {object} dto.MessageResponse
// @Failure      400 {object} dto.ErrorResponse "invalid input"
// @Failure      401 {object} dto.ErrorResponse "authentication required"
// @Failure      403 {object} dto.ErrorResponse "access denied"
// @Failure      404 {object} dto.ErrorResponse "owner request not found"
// @Failure      409 {object} dto.ErrorResponse "owner request is not pending / place owner already set"
// @Failure      500 {object} dto.ErrorResponse "failed to approve owner request"
// @Router       /admin/owner_requests/{id}/approve [post]
func (h *Application) ApproveOwnerRequest(c *gin.Context) {
	if !isPlatformAdmin(c) {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrAccessDenied})
		return
	}

	adminID := c.GetString("user_id")
	if adminID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: response.ErrUnauthorized})
		return
	}

	id := strings.TrimSpace(c.Param("id"))
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	var req dto.ApproveOwnerRequestRequest
	_ = c.ShouldBindJSON(&req) // comment optional

	err := h.OwnerRequestService.Approve(c.Request.Context(), adminID, id, req.Comment)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrOwnerRequestNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: response.ErrOwnerRequestNotFound})
		case errors.Is(err, serviceErrors.ErrOwnerRequestNotPending):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: response.ErrOwnerRequestNotPending})
		case errors.Is(err, serviceErrors.ErrPlaceOwnerAlreadySet):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: response.ErrPlaceOwnerAlreadySet})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrFailedApproveOwnerRequest})
		}
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "approved"})
}

// RejectOwnerRequest godoc
// @Summary      Отклонить заявку владельца (platform admin)
// @Description  Отклоняет заявку (status=rejected)
// @Tags         platform_admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path     string  true  "OwnerRequest ID (uuid)"
// @Param        body  body     dto.RejectOwnerRequestRequest  false  "comment"
// @Success      200 {object} dto.MessageResponse
// @Failure      400 {object} dto.ErrorResponse "invalid input"
// @Failure      401 {object} dto.ErrorResponse "authentication required"
// @Failure      403 {object} dto.ErrorResponse "access denied"
// @Failure      404 {object} dto.ErrorResponse "owner request not found"
// @Failure      409 {object} dto.ErrorResponse "owner request is not pending"
// @Failure      500 {object} dto.ErrorResponse "failed to reject owner request"
// @Router       /admin/owner_requests/{id}/reject [post]
func (h *Application) RejectOwnerRequest(c *gin.Context) {
	if !isPlatformAdmin(c) {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrAccessDenied})
		return
	}

	adminID := c.GetString("user_id")
	if adminID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: response.ErrUnauthorized})
		return
	}

	id := strings.TrimSpace(c.Param("id"))
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: response.ErrInvalidInput})
		return
	}

	var req dto.RejectOwnerRequestRequest
	_ = c.ShouldBindJSON(&req)

	err := h.OwnerRequestService.Reject(c.Request.Context(), adminID, id, req.Comment)
	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrOwnerRequestNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: response.ErrOwnerRequestNotFound})
		case errors.Is(err, serviceErrors.ErrOwnerRequestNotPending):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: response.ErrOwnerRequestNotPending})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrFailedRejectOwnerRequest})
		}
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "rejected"})
}

// ListPendingOwnerRequests godoc
// @Summary      Pending заявки на владельца (platform admin)
// @Tags         platform_admin
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "limit (default 50)"
// @Success      200 {object} dto.ListOwnerRequestsResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /admin/owner_requests/pending [get]
func (h *Application) ListPendingOwnerRequests(c *gin.Context) {
	// если ты оставляешь whitelist-проверку в контроллере:
	if !isPlatformAdmin(c) {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: response.ErrAccessDenied})
		return
	}

	limit := 50
	items, err := h.OwnerRequestService.ListPending(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: response.ErrInternalServerError})
		return
	}

	resp := dto.ListOwnerRequestsResponse{Items: make([]dto.OwnerRequestItem, 0, len(items))}
	for _, it := range items {
		var reviewedAt *string
		if it.ReviewedAt != nil {
			s := it.ReviewedAt.UTC().Format(time.RFC3339)
			reviewedAt = &s
		}
		var reviewedBy *string
		if it.ReviewedBy != nil {
			s := it.ReviewedBy.String()
			reviewedBy = &s
		}

		resp.Items = append(resp.Items, dto.OwnerRequestItem{
			ID:         it.ID.String(),
			UserID:     it.UserID.String(),
			Source:     it.Source,
			SourceID:   it.SourceID,
			Name:       it.Name,
			Status:     string(it.Status),
			CreatedAt:  it.CreatedAt.UTC().Format(time.RFC3339),
			ReviewedAt: reviewedAt,
			ReviewedBy: reviewedBy,
			Comment:    it.Comment,
		})
	}

	c.JSON(http.StatusOK, resp)
}

func isPlatformAdmin(c *gin.Context) bool {
	uid := c.GetString("user_id")
	if uid == "" {
		return false
	}

	raw := os.Getenv("PLATFORM_ADMIN_IDS")
	if raw == "" {
		return false
	}

	for _, id := range strings.Split(raw, ",") {
		if strings.TrimSpace(id) == uid {
			return true
		}
	}
	return false
}
