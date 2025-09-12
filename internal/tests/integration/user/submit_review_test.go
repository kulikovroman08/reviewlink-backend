package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
	s.TS.TruncateAll()

	db := stdlib.OpenDBFromPool(s.TS.DB)
	defer db.Close()

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
	s.Require().NoError(err)
	s.Require().NoError(fixture.Load())

	s.Token = s.TS.Login("bob@example.com", "password123")
}

func (s *SubmitReviewTestSuite) TearDownTest() {
	s.TS.TruncateAll()
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

	s.Require().Equal(http.StatusCreated, rec.Code)
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

	s.Require().Equal(http.StatusBadRequest, rec.Code)
	s.Require().Contains(rec.Body.String(), "invalid token")
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

	s.Require().Equal(http.StatusForbidden, rec.Code)
	s.Require().Contains(rec.Body.String(), "token expired")
}
