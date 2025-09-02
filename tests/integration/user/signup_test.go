package user

import (
	"encoding/json"
	"fmt"
	"github.com/kulikovroman08/reviewlink-backend/tests/integration"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"
	"github.com/stretchr/testify/suite"
)

type SignupTestSuite struct {
	suite.Suite
	TS *integration.TestSetup
}

func TestSignupSuite(t *testing.T) {
	suite.Run(t, new(SignupTestSuite))
}

func (s *SignupTestSuite) SetupTest() {
	s.TS = integration.NewTestSetup()
	s.TS.TruncateUsers()
}

func (s *SignupTestSuite) TestSignupSuccess() {
	body := `{"name": "Alice", "email": "alice@example.com", "password": "123456"}`
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.TS.App.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var resp dto.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.NotEmpty(resp.Token)
}

func (s *SignupTestSuite) TestSignupEmailAlreadyUsed() {
	user := map[string]string{
		"name":     "Alice",
		"email":    "alice@example.com",
		"password": "123456",
	}
	s.createUser(user)

	body := fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s"}`, user["name"], user["email"], user["password"])
	w, resp := s.signupRequest(body)

	s.Equal(http.StatusConflict, w.Code)
	s.Equal("email already in use", resp["error"])
}

func (s *SignupTestSuite) TestSignupMissingName() {
	body := `{"email": "test@example.com", "password": "123456"}`
	w, resp := s.signupRequest(body)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Equal("invalid input", resp["error"])
}

func (s *SignupTestSuite) TestSignupMissingEmail() {
	body := `{"name": "Bob", "password": "123456"}`
	w, resp := s.signupRequest(body)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Equal("invalid input", resp["error"])
}

func (s *SignupTestSuite) TestSignupInvalidEmail() {
	body := `{"name": "Bob", "email": "abc", "password": "123456"}`
	w, resp := s.signupRequest(body)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Equal("invalid input", resp["error"])
}

func (s *SignupTestSuite) TestSignupMissingPassword() {
	body := `{"name": "Charlie", "email": "charlie@example.com"}`
	w, resp := s.signupRequest(body)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Equal("invalid input", resp["error"])
}

func (s *SignupTestSuite) TestSignupShortPassword() {
	body := `{"name": "Dave", "email": "dave@example.com", "password": "123"}`
	w, resp := s.signupRequest(body)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Equal("invalid input", resp["error"])
}

func (s *SignupTestSuite) TestSignupUserNameTooLong() {
	longName := strings.Repeat("a", 300)
	body := fmt.Sprintf(`{"name":"%s","email":"longname@example.com","password":"123456"}`, longName)
	w, resp := s.signupRequest(body)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.Equal("failed to signup", resp["error"])
}

func (s *SignupTestSuite) signupRequest(body string) (*httptest.ResponseRecorder, map[string]string) {
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.TS.App.ServeHTTP(w, req)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	return w, resp
}

func (s *SignupTestSuite) createUser(user map[string]string) {
	body := fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s"}`, user["name"], user["email"], user["password"])
	w, _ := s.signupRequest(body)
	s.Equal(http.StatusOK, w.Code)
}
