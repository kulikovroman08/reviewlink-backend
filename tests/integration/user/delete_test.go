package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kulikovroman08/reviewlink-backend/tests/integration"
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

func (s *DeleteUserTestSuite) SetupTest() {
	s.TS = integration.NewTestSetup()
	s.TS.TruncateUsers()
	s.Token = s.TS.SignupAndLogin("delete@example.com", "password123")
}

func (s *DeleteUserTestSuite) TearDownTest() {
	s.TS.TruncateUsers()
}

func (s *DeleteUserTestSuite) TestDeleteUserSuccess() {
	req := httptest.NewRequest(http.MethodDelete, "/users", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.Token))
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var resp map[string]string
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "user deleted", strings.ToLower(resp["message"]))
}

func (s *DeleteUserTestSuite) TestDeleteUser_Unauthorized() {
	req := httptest.NewRequest(http.MethodDelete, "/users", nil)
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}
