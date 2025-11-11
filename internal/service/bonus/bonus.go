package bonus

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/kulikovroman08/reviewlink-backend/configs"
	"github.com/kulikovroman08/reviewlink-backend/internal/model"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository"
	srvErrors "github.com/kulikovroman08/reviewlink-backend/internal/service/errors"
)

var localRand = rand.New(rand.NewSource(time.Now().UnixNano()))

type bonusService struct {
	userRepo  repository.UserRepository
	bonusRepo repository.BonusRepository
	cfg       *configs.Config
}

func NewBonusService(userRepo repository.UserRepository, bonusRepo repository.BonusRepository, cfg *configs.Config) *bonusService {
	return &bonusService{
		userRepo:  userRepo,
		bonusRepo: bonusRepo,
		cfg:       cfg,
	}
}

func (s *bonusService) RedeemBonus(ctx context.Context, userID, placeID, rewardType string) (*model.BonusReward, error) {
	uuidUser, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	uuidPlace, err := uuid.Parse(placeID)
	if err != nil {
		return nil, fmt.Errorf("invalid place id: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	required := s.cfg.BonusRequiredPts
	if user.Points < required {
		return nil, srvErrors.ErrNotEnoughPoints
	}

	bonus := &model.BonusReward{
		ID:             uuid.New(),
		UserID:         uuidUser,
		PlaceID:        uuidPlace,
		RequiredPoints: required,
		RewardType:     rewardType,
		QRToken:        generateQRToken(),
		IsUsed:         false,
		UsedAt:         nil,
	}

	if err := s.bonusRepo.CreateBonus(ctx, bonus); err != nil {
		return nil, srvErrors.ErrBonusCreateFail
	}

	if err := s.userRepo.RedeemPoints(ctx, userID, required); err != nil {
		return nil, fmt.Errorf("redeem user points: %w", err)
	}

	return bonus, nil
}

func (s *bonusService) GetUserBonuses(ctx context.Context, userID string) ([]model.BonusReward, error) {
	bonuses, err := s.bonusRepo.GetBonusesByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get bonuses by user: %v", err)
	}
	return bonuses, nil
}

func (s *bonusService) ValidateBonus(ctx context.Context, qrToken string) error {
	bonusItem, err := s.bonusRepo.GetByQRToken(ctx, qrToken)
	if err != nil {
		return fmt.Errorf("get bonus by token: %w", err)
	}

	if bonusItem.IsUsed {
		return srvErrors.ErrBonusAlreadyUsed
	}

	if err := s.bonusRepo.MarkBonusUsed(ctx, qrToken); err != nil {
		return fmt.Errorf("mark bonus used: %w", err)
	}

	return nil
}

func generateQRToken() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[localRand.Intn(len(charset))]
	}
	return string(b)
}
