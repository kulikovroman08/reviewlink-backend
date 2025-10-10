package reviewlink

import (
	"bytes"
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

type UpdateReviewTestSuite struct {
	suite.Suite
	TS *integration.TestSetup
}

func TestUpdateReviewSuite(t *testing.T) {
	suite.Run(t, new(UpdateReviewTestSuite))
}

func (s *UpdateReviewTestSuite) SetupSuite() {
	s.TS = integration.NewTestSetup()
}

func (s *UpdateReviewTestSuite) TearDownSuite() {
	s.TS.Close()
}

func (s *UpdateReviewTestSuite) SetupTest() {
	db := stdlib.OpenDBFromPool(s.TS.DB)
	defer db.Close()

	fixture, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			"../fixtures/users.yml",
			"../fixtures/places.yml",
			"../fixtures/get_reviews/reviews.yml",
		),
	)
	require.NoError(s.T(), err)
	require.NoError(s.T(), fixture.Load())
}

func (s *UpdateReviewTestSuite) TestUpdateOwnReviewSuccess() {
	token := s.TS.Login("bob@example.com", "password123")
	require.NotEmpty(s.T(), token)

	body := map[string]any{
		"content": "Обновлённый отзыв — стало ещё лучше!",
		"rating":  4,
	}
	data, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut,
		"/reviews/91c4f4d2-9b0e-4e82-9c5a-9b3a7f7c1a11", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var resp map[string]any
	require.NoError(s.T(), json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(s.T(), "Обновлённый отзыв — стало ещё лучше!", resp["Content"])
	require.Equal(s.T(), float64(4), resp["Rating"])
	require.NotEmpty(s.T(), resp["UpdatedAt"])
}

func (s *UpdateReviewTestSuite) TestUpdateOtherUserReviewForbidden() {
	token := s.TS.Login("john@example.com", "securepass")
	require.NotEmpty(s.T(), token)

	body := map[string]any{
		"content": "Попытка изменить чужой отзыв",
		"rating":  1,
	}
	data, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut,
		"/reviews/91c4f4d2-9b0e-4e82-9c5a-9b3a7f7c1a11", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusForbidden, rec.Code)
}

func (s *UpdateReviewTestSuite) TestUpdateReviewUnauthorized() {
	body := map[string]any{
		"content": "Без токена",
		"rating":  3,
	}
	data, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut,
		"/reviews/91c4f4d2-9b0e-4e82-9c5a-9b3a7f7c1a11", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}
