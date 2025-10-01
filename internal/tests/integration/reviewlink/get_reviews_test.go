package reviewlink

import (
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

type GetReviewsTestSuite struct {
	suite.Suite
	TS *integration.TestSetup
}

func TestGetReviewsSuite(t *testing.T) {
	suite.Run(t, new(GetReviewsTestSuite))
}

func (s *GetReviewsTestSuite) SetupSuite() {
	s.TS = integration.NewTestSetup()
}

func (s *GetReviewsTestSuite) TearDownSuite() {
	s.TS.Close()
}

func (s *GetReviewsTestSuite) SetupTest() {
	db := stdlib.OpenDBFromPool(s.TS.DB)
	defer db.Close()

	fixture, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			"../fixtures/users.yml",
			"../fixtures/places.yml",
			"../fixtures/review_tokens.yml",
			"../fixtures/get_reviews/reviews.yml",
		),
	)
	require.NoError(s.T(), err)
	require.NoError(s.T(), fixture.Load())
}

// 1. Получение всех отзывов
func (s *GetReviewsTestSuite) TestGetAllReviews() {
	const testPlaceID = "a8c52b0c-8f11-4b9c-9c3f-123456789abc"

	req := httptest.NewRequest(http.MethodGet, "/places/"+testPlaceID+"/reviews", nil)
	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var resp []map[string]any
	require.NoError(s.T(), json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(s.T(), resp, 4)
}

func (s *GetReviewsTestSuite) TestFilterByRating5() {
	const (
		testPlaceID = "a8c52b0c-8f11-4b9c-9c3f-123456789abc"
		testRating  = 5
	)

	req := httptest.NewRequest(http.MethodGet, "/places/"+testPlaceID+"/reviews?rating=5", nil)
	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var resp []map[string]any
	require.NoError(s.T(), json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(s.T(), resp, 2)

	for _, r := range resp {
		require.Equal(s.T(), float64(testRating), r["rating"])
	}
}

func (s *GetReviewsTestSuite) TestSortDesc() {
	const testPlaceID = "a8c52b0c-8f11-4b9c-9c3f-123456789abc"

	req := httptest.NewRequest(http.MethodGet, "/places/"+testPlaceID+"/reviews?sort=date_desc", nil)
	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var resp []map[string]any
	require.NoError(s.T(), json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(s.T(), resp, 4)
}
