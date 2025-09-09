package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kulikovroman08/reviewlink-backend/internal/controller/dto"

	"github.com/kulikovroman08/reviewlink-backend/tests/integration"

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

func (s *SignupTestSuite) sendSignupRequest(name, email, password string) (int, map[string]string) {
	body := fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s"}`, name, email, password)
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.TS.App.ServeHTTP(w, req)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	return w.Code, resp
}

func (s *SignupTestSuite) createUser(name, email, password string) {
	code, _ := s.sendSignupRequest(name, email, password)
	s.Equal(http.StatusOK, code)
}

func (s *SignupTestSuite) TestSignupSuccess() {
	code := func() int {
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
		return w.Code
	}()

	s.Equal(http.StatusOK, code)
}

func (s *SignupTestSuite) TestSignupEmailAlreadyUsed() {
	s.createUser("Alice", "alice@example.com", "123456")

	code, resp := s.sendSignupRequest("Alice", "alice@example.com", "123456")
	s.Equal(http.StatusConflict, code)
	s.Equal("email already in use", resp["error"])
}

func (s *SignupTestSuite) TestSignupMissingName() {
	code, resp := s.sendSignupRequest("", "test@example.com", "123456")
	s.Equal(http.StatusBadRequest, code)
	s.Equal("invalid input", resp["error"])
}

func (s *SignupTestSuite) TestSignupInvalidEmail() {
	code, resp := s.sendSignupRequest("Bob", "invalid@@email", "123456")
	s.Equal(http.StatusBadRequest, code)
	s.Equal("invalid input", resp["error"])
}

func (s *SignupTestSuite) TestSignupMissingPassword() {
	code, resp := s.sendSignupRequest("Charlie", "charlie@example.com", "")
	s.Equal(http.StatusBadRequest, code)
	s.Equal("invalid input", resp["error"])
}

func (s *SignupTestSuite) TestSignupShortPassword() {
	code, resp := s.sendSignupRequest("Dave", "dave@example.com", "123")
	s.Equal(http.StatusBadRequest, code)
	s.Equal("invalid input", resp["error"])
}

func (s *SignupTestSuite) TestSignupUserNameTooLong() {
	longName := strings.Repeat("a", 300)
	code, resp := s.sendSignupRequest(longName, "long@example.com", "123456")

	s.Equal(http.StatusBadRequest, code)
	s.Equal("invalid input", resp["error"])
}
