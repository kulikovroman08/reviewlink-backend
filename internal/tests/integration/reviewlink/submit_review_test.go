package reviewlink

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/kulikovroman08/reviewlink-backend/internal/tests/integration"

	"github.com/stretchr/testify/suite"
)

type SubmitReviewTestSuite struct {
	suite.Suite
	TS    *integration.TestSetup
	Token string
}

func TestSubmitReviewSuite(t *testing.T) {
	suite.Run(t, new(SubmitReviewTestSuite))
}

func (s *SubmitReviewTestSuite) SetupSuite() {
	s.TS = integration.NewTestSetup()
}

func (s *SubmitReviewTestSuite) TearDownSuite() {
	s.TS.Close()
}

func (s *SubmitReviewTestSuite) SetupTest() {
	db := stdlib.OpenDBFromPool(s.TS.DB)
	defer func() {
		if err := db.Close(); err != nil {
			s.T().Logf("failed to close db: %v", err)
		}
	}()

	fixture, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			"../fixtures/users.yml",
			"../fixtures/places.yml",
			"../fixtures/review_tokens.yml",
			"../fixtures/reviews.yml",
		),
	)

	require.NoError(s.T(), err, "init fixtures failed")
	require.NoError(s.T(), fixture.Load(), "load fixtures failed")

	s.Token = s.TS.Login("bob@example.com", "password123")
}

func (s *SubmitReviewTestSuite) TestSubmitReviewSuccess() {
	const (
		testPlaceID = "a8c52b0c-8f11-4b9c-9c3f-123456789abc"
		testToken   = "VALIDTOKEN123"
		testContent = "Очень понравилось!"
		testRating  = 5
	)

	payload := map[string]any{
		"rating":   testRating,
		"content":  testContent,
		"place_id": testPlaceID,
		"token":    testToken,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusCreated, rec.Code)
}

func (s *SubmitReviewTestSuite) TestSubmitReviewAlreadyToday() {
	const (
		testPlaceID = "a8c52b0c-8f11-4b9c-9c3f-123456789abc"
		testToken   = "USEDTODAY999"
		testContent = "Повторный отзыв"
		testRating  = 4
	)

	s.Token = s.TS.Login("john@example.com", "securepass")

	payload := map[string]any{
		"rating":   testRating,
		"content":  testContent,
		"place_id": testPlaceID,
		"token":    testToken,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusBadRequest, rec.Code)
	require.Contains(s.T(), rec.Body.String(), "invalid credentials")
}

func (s *SubmitReviewTestSuite) TestSubmitReviewExpiredToken() {
	const (
		testPlaceID = "a8c52b0c-8f11-4b9c-9c3f-123456789abc"
		testToken   = "EXPIRED00001"
		testContent = "Старый токен"
		testRating  = 3
	)

	payload := map[string]any{
		"rating":   testRating,
		"content":  testContent,
		"place_id": testPlaceID,
		"token":    testToken,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusForbidden, rec.Code)
	require.Contains(s.T(), rec.Body.String(), "token expired")
}
