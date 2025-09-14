package reviewlink

import (
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

type GetUserTestSuite struct {
	suite.Suite
	TS    *integration.TestSetup
	Token string
}

func TestGetUserSuite(t *testing.T) {
	suite.Run(t, new(GetUserTestSuite))
}

func (s *GetUserTestSuite) SetupSuite() {
	s.TS = integration.NewTestSetup()
}

func (s *GetUserTestSuite) TearDownSuite() {
	s.TS.Close()
}

func (s *GetUserTestSuite) SetupTest() {
	s.TS.TruncateAll()

	db := stdlib.OpenDBFromPool(s.TS.DB)
	defer db.Close()

	fixture, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			"../fixtures/users.yml",
		),
	)
	require.NoError(s.T(), err, "init fixtures failed")
	require.NoError(s.T(), fixture.Load(), "load fixtures failed")

	s.Token = s.TS.Login("john@example.com", "securepass")
}

func (s *GetUserTestSuite) TearDownTest() {
	s.TS.TruncateAll()
}

func (s *GetUserTestSuite) TestGetUserSuccess() {
	const (
		expectedEmail = "john@example.com"
	)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", "Bearer "+s.Token)
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var resp map[string]any
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedEmail, resp["email"])
}

func (s *GetUserTestSuite) TestGetUserUnauthorizedMissingToken() {
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *GetUserTestSuite) TestGetUserUnauthorizedInvalidToken() {
	const (
		invalidToken = "invalid-token"
	)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", "Bearer "+invalidToken)
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}
