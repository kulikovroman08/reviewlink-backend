package reviewlink

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/kulikovroman08/reviewlink-backend/internal/tests/integration"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DeleteUserTestSuite struct {
	suite.Suite
	TS    *integration.TestSetup
	Token string
}

func TestDeleteUserSuite(t *testing.T) {
	suite.Run(t, new(DeleteUserTestSuite))
}

func (s *DeleteUserTestSuite) SetupSuite() {
	s.TS = integration.NewTestSetup()
}

func (s *DeleteUserTestSuite) TearDownSuite() {
	s.TS.Close()
}

func (s *DeleteUserTestSuite) SetupTest() {
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
	s.Require().NoError(err)
	s.Require().NoError(fixture.Load())

	s.Token = s.TS.Login("john@example.com", "securepass")
}

func (s *DeleteUserTestSuite) TearDownTest() {
	s.TS.TruncateAll()
}

func (s *DeleteUserTestSuite) TestDeleteUserSuccess() {
	const (
		expectedMessage = "user deleted"
	)

	req := httptest.NewRequest(http.MethodDelete, "/users", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.Token))
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var resp map[string]string
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedMessage, strings.ToLower(resp["message"]))
}

func (s *DeleteUserTestSuite) TestDeleteUser_Unauthorized() {
	req := httptest.NewRequest(http.MethodDelete, "/users", nil)
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}
