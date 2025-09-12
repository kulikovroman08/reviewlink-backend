package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kulikovroman08/reviewlink-backend/internal/tests/integration"

	"github.com/stretchr/testify/suite"
)

type SignupTestSuite struct {
	suite.Suite
	TS *integration.TestSetup
}

func TestSignupSuite(t *testing.T) {
	suite.Run(t, new(SignupTestSuite))
}

func (s *SignupTestSuite) SetupSuite() {
	s.TS = integration.NewTestSetup()
}

func (s *SignupTestSuite) TearDownSuite() {
	s.TS.Close()
}

func (s *SignupTestSuite) SetupTest() {
	s.TS.TruncateAll()
}

func (s *SignupTestSuite) TearDownTest() {
	s.TS.TruncateAll()
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
	s.Require().Equal(http.StatusOK, code)
}

func (s *SignupTestSuite) assertErrorResponse(code int, resp map[string]string, expectedCode int, expectedErr string) {
	s.Equal(expectedCode, code)
	s.Equal(expectedErr, resp["error"])
}

func (s *SignupTestSuite) TestSignupSuccess() {
	code, resp := s.sendSignupRequest("Alice", "alice@example.com", "123456")
	s.Equal(http.StatusOK, code)
	s.NotEmpty(resp["token"])
}

func (s *SignupTestSuite) TestSignupEmailAlreadyUsed() {
	s.createUser("Alice", "alice@example.com", "123456")

	code, resp := s.sendSignupRequest("Alice", "alice@example.com", "123456")
	s.assertErrorResponse(code, resp, http.StatusConflict, "email already in use")
}

func (s *SignupTestSuite) TestSignupMissingName() {
	code, resp := s.sendSignupRequest("", "test@example.com", "123456")
	s.assertErrorResponse(code, resp, http.StatusBadRequest, "invalid input")
}

func (s *SignupTestSuite) TestSignupInvalidEmail() {
	code, resp := s.sendSignupRequest("Bob", "invalid@@email", "123456")
	s.assertErrorResponse(code, resp, http.StatusBadRequest, "invalid input")
}

func (s *SignupTestSuite) TestSignupMissingPassword() {
	code, resp := s.sendSignupRequest("Charlie", "charlie@example.com", "")
	s.assertErrorResponse(code, resp, http.StatusBadRequest, "invalid input")
}

func (s *SignupTestSuite) TestSignupShortPassword() {
	code, resp := s.sendSignupRequest("Dave", "dave@example.com", "123")
	s.assertErrorResponse(code, resp, http.StatusBadRequest, "invalid input")
}

func (s *SignupTestSuite) TestSignupUserNameTooLong() {
	longName := strings.Repeat("a", 300)
	code, resp := s.sendSignupRequest(longName, "long@example.com", "123456")
	s.assertErrorResponse(code, resp, http.StatusBadRequest, "invalid input")
}
