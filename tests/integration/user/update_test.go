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
	s.TS.TruncateUsers()
	s.Token = s.TS.SignupAndLogin("update@example.com", "password123")
}

func (s *UpdateUserTestSuite) TestUpdateUserNameSuccess() {
	body := map[string]string{
		"name": "New Name",
	}
	data, _ := json.Marshal(body)

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
	body := map[string]string{
		"email": "not-an-email",
	}
	data, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewReader(data))
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusBadRequest, rec.Code)
}
