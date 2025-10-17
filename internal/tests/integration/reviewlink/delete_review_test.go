package reviewlink

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/kulikovroman08/reviewlink-backend/internal/tests/integration"
)

type DeleteReviewTestSuite struct {
	suite.Suite
	TS    *integration.TestSetup
	Token string
}

func TestDeleteReviewSuite(t *testing.T) {
	suite.Run(t, new(DeleteReviewTestSuite))
}

func (s *DeleteReviewTestSuite) SetupSuite() {
	s.TS = integration.NewTestSetup()
}

func (s *DeleteReviewTestSuite) TearDownSuite() {
	s.TS.Close()
}

func (s *DeleteReviewTestSuite) SetupTest() {
	db := stdlib.OpenDBFromPool(s.TS.DB)
	defer db.Close()

	fixture, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			"../fixtures/users.yml",
			"../fixtures/places.yml",
			"../fixtures/review_tokens.yml",
			"../fixtures/delete_reviews/reviews.yml",
		),
	)
	require.NoError(s.T(), err)
	require.NoError(s.T(), fixture.Load())

	s.Token = s.TS.Login("bob@example.com", "password123")
}

func (s *DeleteReviewTestSuite) TestDeleteReviewSuccess() {
	const (
		reviewID   = "01eb078d-a93b-410e-9c83-76251fb04f07"
		successMsg = `{"message":"review deleted successfully"}`
	)

	req := httptest.NewRequest(http.MethodDelete, "/reviews/"+reviewID, nil)
	req.Header.Set("Authorization", "Bearer "+s.Token)
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)
	require.Equal(s.T(), successMsg, rec.Body.String())
}

func (s *DeleteReviewTestSuite) TestDeleteReviewForbidden() {
	const (
		reviewID = "01eb078d-a93b-410e-9c83-76251fb04f07"
		errorMsg = `{"error":"review not found"}`
	)

	token := s.TS.Login("john@example.com", "securepass")

	req := httptest.NewRequest(http.MethodDelete, "/reviews/"+reviewID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusForbidden, rec.Code)
	require.Equal(s.T(), errorMsg, rec.Body.String())
}
