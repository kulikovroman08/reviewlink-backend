package unit

import (
	"context"
	"errors"
	"github.com/kulikovroman08/reviewlink-backend/internal/tests/mocks"
	"go.uber.org/mock/gomock"
	"testing"

	"github.com/kulikovroman08/reviewlink-backend/internal/service/token"
	"github.com/stretchr/testify/require"
)

func TestSuccessGenerateTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTokenRepository(ctrl)
	svc := token.NewTokenService(mockRepo)

	ctx := context.Background()
	placeID := "a8c52b0c-8f11-4b9c-9c3f-123456789abc"
	count := 3

	mockRepo.EXPECT().
		CreateTokens(ctx, gomock.Any()).
		Return(nil)

	result, err := svc.GenerateTokens(ctx, placeID, count)

	require.NoError(t, err)
	require.Len(t, result.Tokens, count)
}

func TestFailGenerateTokens(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTokenRepository(ctrl)
	svc := token.NewTokenService(mockRepo)

	ctx := context.Background()
	placeID := "a8c52b0c-8f11-4b9c-9c3f-123456789abc"
	count := 2

	mockRepo.EXPECT().
		CreateTokens(ctx, gomock.Any()).
		Return(errors.New("db error"))

	result, err := svc.GenerateTokens(ctx, placeID, count)

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "db error")
}

func TestFailGenerateTokensInvalidPlaceID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTokenRepository(ctrl)
	svc := token.NewTokenService(mockRepo)

	ctx := context.Background()
	invalidPlaceID := "not-a-uuid"

	result, err := svc.GenerateTokens(ctx, invalidPlaceID, 1)

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "invalid place_id")
}
