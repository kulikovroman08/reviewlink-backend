package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kulikovroman08/reviewlink-backend/tests/integration"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetUserTestSuite struct {
	suite.Suite
	TS *integration.TestSetup
}

func TestGetUserSuite(t *testing.T) {
	suite.Run(t, new(GetUserTestSuite))
}

func (s *GetUserTestSuite) SetupTest() {
	s.TS = integration.NewTestSetup()
	s.TS.TruncateUsers()
}

func (s *GetUserTestSuite) TestGetUserSuccess() {
	token := s.TS.SignupAndLogin("john@example.com", "securepass")

	req := httptest.NewRequest("GET", "/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var resp map[string]any
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "john@example.com", resp["email"])
}

func (s *GetUserTestSuite) TestGetUserUnauthorizedMissingToken() {
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *GetUserTestSuite) TestGetUserUnauthorizedInvalidToken() {
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

//
//func (s *GetUserTestSuite) TestGetUserNotFound() {
//	req := httptest.NewRequest(http.MethodGet, "/users", nil)
//	req.Header.Set("X-Test-User-ID", "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
//	resp := httptest.NewRecorder()
//
//	s.TS.App.ServeHTTP(resp, req)
//
//	s.Equal(http.StatusNotFound, resp.Code)
//
//	var res map[string]interface{}
//	err := json.Unmarshal(resp.Body.Bytes(), &res)
//	s.Require().NoError(err)
//	s.Contains(res["error"], "not found")
//}
//
//func (s *GetUserTestSuite) TestGetUserDeletedUser() {
//	req := httptest.NewRequest(http.MethodGet, "/users", nil)
//	req.Header.Set("X-Test-User-ID", "1c1f91a2-1234-4e22-aaaa-111111111111") // Deleted User
//	resp := httptest.NewRecorder()
//
//	s.TS.App.ServeHTTP(resp, req)
//
//	s.Equal(http.StatusNotFound, resp.Code)
//
//	var res map[string]interface{}
//	err := json.Unmarshal(resp.Body.Bytes(), &res)
//	s.Require().NoError(err)
//	s.Equal("user not found", res["error"])
//}
//
//func (s *GetUserTestSuite) TestGetUser_InvalidUUIDFormat() {
//	req := httptest.NewRequest(http.MethodGet, "/users", nil)
//	req.Header.Set("X-Test-User-ID", "not-a-valid-uuid")
//	resp := httptest.NewRecorder()
//
//	s.TS.App.ServeHTTP(resp, req)
//
//	s.Equal(http.StatusInternalServerError, resp.Code)
//
//	var res map[string]interface{}
//	err := json.Unmarshal(resp.Body.Bytes(), &res)
//	s.Require().NoError(err)
//	s.Contains(res["error"], "failed to get user") // контроллер отлавливает 22P02
//}
//
//func (s *GetUserTestSuite) TestGetUserMissingUserID() {
//	req := httptest.NewRequest(http.MethodGet, "/users", nil) // no X-Test-User-ID
//	resp := httptest.NewRecorder()
//
//	s.TS.App.ServeHTTP(resp, req)
//
//	s.Equal(http.StatusUnauthorized, resp.Code)
//
//	var res map[string]interface{}
//	err := json.Unmarshal(resp.Body.Bytes(), &res)
//	s.Require().NoError(err)
//	s.Equal("Authorization header required", res["error"])
//}
