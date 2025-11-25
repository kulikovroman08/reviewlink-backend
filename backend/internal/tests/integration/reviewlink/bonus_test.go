package reviewlink

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/kulikovroman08/reviewlink-backend/internal/tests/integration"
)

type BonusTestSuite struct {
	suite.Suite
	TS *integration.TestSetup
}

func TestBonusSuite(t *testing.T) {
	suite.Run(t, new(BonusTestSuite))
}

func (s *BonusTestSuite) SetupSuite() {
	s.TS = integration.NewTestSetup()
}

func (s *BonusTestSuite) TearDownSuite() {
	s.TS.Close()
}

func (s *BonusTestSuite) SetupTest() {
	db := stdlib.OpenDBFromPool(s.TS.DB)
	defer db.Close()

	fixture, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			"../fixtures/users.yml",
			"../fixtures/places.yml",
			"../fixtures/bonuses/bonus_rewards.yml",
		),
	)
	require.NoError(s.T(), err)
	require.NoError(s.T(), fixture.Load())
}

func (s *BonusTestSuite) TestGetUserBonuses() {
	token := s.TS.Login("john@example.com", "securepass")

	req := httptest.NewRequest(http.MethodGet, "/bonuses", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var resp []map[string]any
	require.NoError(s.T(), json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(s.T(), resp, 1)
	require.Equal(s.T(), "free_coffee", resp[0]["reward_type"])
	require.Equal(s.T(), false, resp[0]["is_used"])
}

func (s *BonusTestSuite) TestRedeemBonusSuccess() {
	token := s.TS.Login("john@example.com", "securepass")

	body := map[string]any{
		"place_id":    "a8c52b0c-8f11-4b9c-9c3f-123456789abc",
		"reward_type": "free_meal",
	}
	data, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/bonuses/redeem", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	// ✅ ожидаем успех, 201
	require.Equal(s.T(), http.StatusCreated, rec.Code)

	var resp map[string]any
	require.NoError(s.T(), json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(s.T(), "free_meal", resp["reward_type"])
	require.NotEmpty(s.T(), resp["qr_token"])
}

func (s *BonusTestSuite) TestRedeemBonusNotEnoughPoints() {
	token := s.TS.Login("bob@example.com", "password123")

	_, err := s.TS.DB.Exec(context.Background(),
		"UPDATE users SET points = 10 WHERE email = $1", "bob@example.com")
	require.NoError(s.T(), err)

	body := map[string]any{
		"place_id":    "a8c52b0c-8f11-4b9c-9c3f-123456789abc",
		"reward_type": "free_coffee",
	}
	data, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/bonuses/redeem", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusConflict, rec.Code)

	var resp map[string]any
	require.NoError(s.T(), json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(s.T(), "not enough points", resp["error"])
}

func (s *BonusTestSuite) TestValidateBonusSuccess() {
	token := s.TS.Login("admin@example.com", "securepass")

	body := map[string]any{
		"qr_token": "bonus123qr",
	}
	data, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/bonuses/validate", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var resp map[string]any
	require.NoError(s.T(), json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(s.T(), "bonus redeemed", resp["status"])
}
