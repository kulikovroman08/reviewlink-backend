package reviewlink

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

type UpdateUserTestSuite struct {
	suite.Suite
	TS    *integration.TestSetup
	Token string
}

func TestUpdateUserSuite(t *testing.T) {
	suite.Run(t, new(UpdateUserTestSuite))
}

func (s *UpdateUserTestSuite) SetupTest() {
	s.TS = integration.NewTestSetup()

	db := stdlib.OpenDBFromPool(s.TS.DB)
	defer func() {
		if err := db.Close(); err != nil {
			s.T().Logf("failed to close db: %v", err)
		}
	}()

	fixture, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			"../fixtures/users.yml",
		),
	)
	require.NoError(s.T(), err)
	require.NoError(s.T(), fixture.Load())

	s.Token = s.TS.Login("update@example.com", "password123")
}

func (s *UpdateUserTestSuite) TestUpdateUserNameSuccess() {
	const (
		newName = "New Name"
	)

	body := map[string]string{
		"name": newName,
	}
	data, err := json.Marshal(body)
	require.NoError(s.T(), err)

	req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusOK, rec.Code)
}

func (s *UpdateUserTestSuite) TestUpdateUserNoFieldsProvided() {
	data := []byte(`{}`)

	req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusBadRequest, rec.Code)
}

func (s *UpdateUserTestSuite) TestUpdateUserInvalidEmail() {
	const (
		invalidEmail = "not-an-email"
	)

	body := map[string]string{
		"email": invalidEmail,
	}
	data, err := json.Marshal(body)
	require.NoError(s.T(), err)

	req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusBadRequest, rec.Code)
}
