package user

import (
	"bytes"
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

type LoginTestSuite struct {
	suite.Suite
	TS *integration.TestSetup
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}

func (s *LoginTestSuite) SetupSuite() {
	s.TS = integration.NewTestSetup()
}

func (s *LoginTestSuite) TearDownSuite() {
	s.TS.Close()
}

func (s *LoginTestSuite) SetupTest() {
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
}

func (s *LoginTestSuite) TearDownTest() {
	s.TS.TruncateAll()
}

func (s *LoginTestSuite) TestLoginSuccess() {
	const (
		inputEmail    = "bob@example.com"
		inputPassword = "password123"
	)

	loginPayload := map[string]string{
		"email":    inputEmail,
		"password": inputPassword,
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
	const (
		inputEmail    = "bob@example.com"
		wrongPassword = "wrongpassword"
	)

	payload := map[string]string{
		"email":    inputEmail,
		"password": wrongPassword,
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
	const (
		nonExistentEmail = "ghost@example.com"
		anyPassword      = "password123"
	)

	payload := map[string]string{
		"email":    nonExistentEmail,
		"password": anyPassword,
	}
	body, err := json.Marshal(payload)
	require.NoError(s.T(), err)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}
