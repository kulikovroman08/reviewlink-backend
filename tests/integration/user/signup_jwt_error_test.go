package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/kulikovroman08/reviewlink-backend/tests/integration"
	"github.com/stretchr/testify/suite"
)

type SignupJWTErrorSuite struct {
	suite.Suite
	TS       *integration.TestSetup
	original string
}

func TestSignupJWTErrorSuite(t *testing.T) {
	suite.Run(t, new(SignupJWTErrorSuite))
}

func (s *SignupJWTErrorSuite) SetupTest() {
	s.original = os.Getenv("JWT_SECRET")

	_ = os.Setenv("JWT_SECRET", "")

	s.TS = integration.NewTestSetup()
	s.TS.TruncateUsers()
}

func (s *SignupJWTErrorSuite) TearDownTest() {
	_ = os.Setenv("JWT_SECRET", s.original)
}

func (s *SignupJWTErrorSuite) TestJWTGenerationFails() {
	body := `{"name": "JWTFail", "email": "jwtfail@example.com", "password": "123456"}`
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.TS.App.ServeHTTP(w, req)

	s.Equal(http.StatusInternalServerError, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	s.NoError(err)
	s.Equal("failed to signup", resp["error"])
}
