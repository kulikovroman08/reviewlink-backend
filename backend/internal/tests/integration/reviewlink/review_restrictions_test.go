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
	"github.com/kulikovroman08/reviewlink-backend/internal/tests/integration"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ReviewRestrictionsTestSuite struct {
	suite.Suite
	TS *integration.TestSetup
}

func TestReviewRestrictionsSuite(t *testing.T) {
	suite.Run(t, new(ReviewRestrictionsTestSuite))
}

func (s *ReviewRestrictionsTestSuite) SetupSuite() {
	s.TS = integration.NewTestSetup()
}

func (s *ReviewRestrictionsTestSuite) TearDownSuite() {
	s.TS.Close()
}

func (s *ReviewRestrictionsTestSuite) SetupTest() {
	db := stdlib.OpenDBFromPool(s.TS.DB)
	defer db.Close()

	fixture, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			"../fixtures/restrictions/users.yml",
			"../fixtures/restrictions/places.yml",
			"../fixtures/restrictions/review_tokens.yml",
			"../fixtures/restrictions/user_restrictions.yml",
		),
	)
	require.NoError(s.T(), err)
	require.NoError(s.T(), fixture.Load())
	_, _ = s.TS.DB.Exec(context.Background(), "DELETE FROM reviews")
}

func (s *ReviewRestrictionsTestSuite) TestActiveRestriction_NoPointsAwarded() {
	token := s.TS.Login("limited_user@example.com", "password123")

	payload := map[string]any{
		"rating":   5,
		"content":  "Хорошее место!",
		"place_id": "d8c52b0c-8f11-4b9c-9c3f-123456789bcd",
		"token":    "LIMITEDTOKEN999",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusCreated, rec.Code, "Review should be created but no points")
}
func (s *ReviewRestrictionsTestSuite) TestNoRestriction_PointsAwarded() {
	token := s.TS.Login("normal_user@example.com", "password123")

	payload := map[string]any{
		"rating":   5,
		"content":  "Отличный сервис!",
		"place_id": "c8c52b0c-8f11-4b9c-9c3f-123456789abc",
		"token":    "FRESHTOKEN123",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusCreated, rec.Code, "Review should be created with points")
}
