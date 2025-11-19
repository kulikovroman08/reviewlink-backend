package reviewlink

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/kulikovroman08/reviewlink-backend/internal/tests/integration"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type LeaderboardTestSuite struct {
	suite.Suite
	TS *integration.TestSetup
}

func TestLeaderboardSuite(t *testing.T) {
	suite.Run(t, new(LeaderboardTestSuite))
}

func (s *LeaderboardTestSuite) SetupSuite() {
	s.TS = integration.NewTestSetup()
}

func (s *LeaderboardTestSuite) TearDownSuite() {
	s.TS.Close()
}

func (s *LeaderboardTestSuite) SetupTest() {
	db := stdlib.OpenDBFromPool(s.TS.DB)
	defer db.Close()

	fixture, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			"../fixtures/users.yml",
			"../fixtures/places.yml",
			"../fixtures/admin/reviews.yml",
		),
	)
	require.NoError(s.T(), err)
	require.NoError(s.T(), fixture.Load())
}

func (s *LeaderboardTestSuite) TestSuccessGetUserLeaderboard() {
	req := httptest.NewRequest("GET", "/leaderboard/users?sort_by=rating&limit=5", nil)
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)
	require.Equal(s.T(), http.StatusOK, rec.Code)

	var resp []dto.LeaderboardEntry
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(s.T(), err)

	require.Len(s.T(), resp, 3)

	require.Equal(s.T(), "John", resp[0].Name)
	require.Equal(s.T(), 5.0, resp[0].AvgRating)

	require.Equal(s.T(), "Admin", resp[1].Name)
	require.Equal(s.T(), 4.0, resp[1].AvgRating)

	require.Equal(s.T(), "Bob", resp[2].Name)
	require.Equal(s.T(), 3.0, resp[2].AvgRating)
}

func (s *LeaderboardTestSuite) Test_GetUserLeaderboard_InvalidSortBy() {
	req := httptest.NewRequest("GET", "/leaderboard/users?sort_by=invalid", nil)
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)
	require.Equal(s.T(), http.StatusOK, rec.Code)

	var resp []dto.LeaderboardEntry
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(s.T(), err)

	require.Len(s.T(), resp, 3)

	require.Equal(s.T(), "John", resp[0].Name)
	require.Equal(s.T(), "Admin", resp[1].Name)
	require.Equal(s.T(), "Bob", resp[2].Name)
}

func (s *LeaderboardTestSuite) Test_GetUserLeaderboard_EmptyResult() {
	req := httptest.NewRequest("GET", "/leaderboard/users?min_rating=10", nil)
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)
	require.Equal(s.T(), http.StatusOK, rec.Code)

	var resp []dto.LeaderboardEntry
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(s.T(), err)

	require.Empty(s.T(), resp)
}
