package reviewlink

import (
	"encoding/json"
	"fmt"
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

type AdminStatsTestSuite struct {
	suite.Suite
	TS    *integration.TestSetup
	Token string
}

func TestAdminStatsSuite(t *testing.T) {
	suite.Run(t, new(AdminStatsTestSuite))
}

func (s *AdminStatsTestSuite) SetupSuite() {
	s.TS = integration.NewTestSetup()
}

func (s *AdminStatsTestSuite) TearDownSuite() {
	s.TS.Close()
}

func (s *AdminStatsTestSuite) SetupTest() {
	db := stdlib.OpenDBFromPool(s.TS.DB)
	defer db.Close()

	fixture, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			"../fixtures/users.yml",
			"../fixtures/admin/reviews.yml",
			"../fixtures/admin/bonus_rewards.yml",
		),
	)
	require.NoError(s.T(), err, "init fixtures failed")
	require.NoError(s.T(), fixture.Load(), "load fixtures failed")

	s.Token = s.TS.Login("admin@example.com", "securepass")
}

func (s *AdminStatsTestSuite) TestAdminStats_Success() {
	req := httptest.NewRequest(http.MethodGet, "/admin/stats", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.Token))
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var resp dto.AdminStatsResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(s.T(), err)

	require.Equal(s.T(), 4, resp.TotalUsers)
	require.Equal(s.T(), 3, resp.TotalReviews)
	require.Equal(s.T(), 4.0, resp.AverageRating)
	require.Equal(s.T(), 2, resp.TotalBonuses)
}

func (s *AdminStatsTestSuite) TestAdminStats_Forbidden() {
	token := s.TS.Login("bob@example.com", "password123")

	req := httptest.NewRequest(http.MethodGet, "/admin/stats", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusForbidden, rec.Code)
}
