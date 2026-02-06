package owner_request

import (
	"context"
	"fmt"

	"github.com/kulikovroman08/reviewlink-backend/internal/service/errors"

	"github.com/google/uuid"

	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository"
)

type ownerRequestService struct {
	ownerRequestRepo repository.OwnerRequestRepository
	placeRepo        repository.PlaceRepository
	userRepo         repository.UserRepository
}

func NewOwnerRequestService(
	ownerRequestRepo repository.OwnerRequestRepository,
	placeRepo repository.PlaceRepository,
	userRepo repository.UserRepository,
) *ownerRequestService {
	return &ownerRequestService{
		ownerRequestRepo: ownerRequestRepo,
		placeRepo:        placeRepo,
		userRepo:         userRepo,
	}
}

func (s *ownerRequestService) Create(ctx context.Context, userID string, source string, sourceID string, name string) (string, error) {
	if userID == "" || source == "" || sourceID == "" || name == "" {
		return "", errors.ErrInvalidOwnerRequestData
	}

	uID, err := uuid.Parse(userID)
	if err != nil {
		return "", fmt.Errorf("invalid user_id: %w", err)
	}

	req := model.OwnerRequest{
		ID:       uuid.New(),
		UserID:   uID,
		Source:   source,
		SourceID: sourceID,
		Name:     name,
		Status:   model.OwnerRequestPending,
	}

	if err := s.ownerRequestRepo.Create(ctx, &req); err != nil {
		return "", fmt.Errorf("failed to create owner request: %w", err)
	}

	return req.ID.String(), nil
}

func (s *ownerRequestService) ListPending(ctx context.Context, limit int) ([]model.OwnerRequest, error) {
	items, err := s.ownerRequestRepo.ListPending(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list pending owner requests: %w", err)
	}
	return items, nil
}

func (s *ownerRequestService) Approve(ctx context.Context, adminID string, requestID string, comment *string) error {
	if adminID == "" {
		return fmt.Errorf("admin_id is required")
	}
	if requestID == "" {
		return fmt.Errorf("request_id is required")
	}
	if _, err := uuid.Parse(adminID); err != nil {
		return fmt.Errorf("invalid admin_id: %w", err)
	}

	req, err := s.ownerRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("failed to get owner request: %w", err)
	}
	if req == nil {
		return errors.ErrOwnerRequestNotFound
	}
	if req.Status != model.OwnerRequestPending {
		return errors.ErrOwnerRequestNotPending
	}

	placeID, err := s.placeRepo.EnsureFromPublic(ctx, req.Source, req.SourceID, req.Name)
	if err != nil {
		return fmt.Errorf("failed to ensure place from public: %w", err)
	}

	if err := s.placeRepo.SetOwner(ctx, placeID, req.UserID.String()); err != nil {
		return errors.ErrPlaceOwnerAlreadySet
	}

	if err := s.ownerRequestRepo.Approve(ctx, requestID, adminID, comment); err != nil {
		return fmt.Errorf("failed to approve owner request: %w", err)
	}

	if err := s.userRepo.SetRole(ctx, req.UserID.String(), "admin"); err != nil {
		return fmt.Errorf("failed to promote user role to admin: %w", err)
	}

	return nil
}

func (s *ownerRequestService) Reject(ctx context.Context, adminID string, requestID string, comment *string) error {
	if adminID == "" {
		return fmt.Errorf("admin_id is required")
	}
	if requestID == "" {
		return fmt.Errorf("request_id is required")
	}
	if _, err := uuid.Parse(adminID); err != nil {
		return fmt.Errorf("invalid admin_id: %w", err)
	}

	req, err := s.ownerRequestRepo.GetByID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("failed to get owner request: %w", err)
	}
	if req == nil {
		return errors.ErrOwnerRequestNotFound
	}
	if req.Status != model.OwnerRequestPending {
		return errors.ErrOwnerRequestNotPending
	}

	if err := s.ownerRequestRepo.Reject(ctx, requestID, adminID, comment); err != nil {
		return fmt.Errorf("failed to reject owner request: %w", err)
	}

	return nil
}
