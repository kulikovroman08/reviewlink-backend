package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kulikovroman08/reviewlink-backend/tests/integration"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type LoginTestSuite struct {
	suite.Suite
	TS *integration.TestSetup
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}

func (s *LoginTestSuite) SetupTest() {
	s.TS = integration.NewTestSetup()
	s.TS.TruncateUsers()

	signupPayload := map[string]string{
		"name":     "Alice",
		"email":    "alice@example.com",
		"password": "password123",
	}
	body, err := json.Marshal(signupPayload)
	require.NoError(s.T(), err)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)
	require.Equal(s.T(), http.StatusOK, rec.Code)
}

func (s *LoginTestSuite) TearDownTest() {
	s.TS.TruncateUsers()
}

func (s *LoginTestSuite) TestLoginSuccess() {
	loginPayload := map[string]string{
		"email":    "alice@example.com",
		"password": "password123",
	}
	body, err := json.Marshal(loginPayload)
	require.NoError(s.T(), err)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)

	var resp map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resp["token"])
}

func (s *LoginTestSuite) TestLoginWrongPassword() {
	payload := map[string]string{
		"email":    "alice@example.com",
		"password": "wrongpassword",
	}
	body, err := json.Marshal(payload)
	require.NoError(s.T(), err)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}

func (s *LoginTestSuite) TestLoginNonExistentUser() {
	payload := map[string]string{
		"email":    "ghost@example.com",
		"password": "password123",
	}
	body, err := json.Marshal(payload)
	require.NoError(s.T(), err)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}
