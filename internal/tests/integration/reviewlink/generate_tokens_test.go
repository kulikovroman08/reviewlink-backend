package reviewlink

//
//import (
//	"bytes"
//	"encoding/json"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/go-testfixtures/testfixtures/v3"
//	"github.com/jackc/pgx/v5/stdlib"
//	"github.com/kulikovroman08/reviewlink-backend/internal/tests/integration"
//
//	"github.com/stretchr/testify/suite"
//)
//
//type GenerateTokensTestSuite struct {
//	suite.Suite
//	TS    *integration.TestSetup
//	Token string
//}
//
//func TestGenerateTokensSuite(t *testing.T) {
//	suite.Run(t, new(GenerateTokensTestSuite))
//}
//
//func (s *GenerateTokensTestSuite) SetupSuite() {
//	s.TS = integration.NewTestSetup()
//}
//
//func (s *GenerateTokensTestSuite) TearDownSuite() {
//	s.TS.Close()
//}
//
//func (s *GenerateTokensTestSuite) SetupTest() {
//	s.TS.TruncateAll()
//
//	db := stdlib.OpenDBFromPool(s.TS.DB)
//	defer db.Close()
//
//	fixture, err := testfixtures.New(
//		testfixtures.Database(db),
//		testfixtures.Dialect("postgres"),
//		testfixtures.Files(
//			"../fixtures/users.yml",
//			"../fixtures/places.yml",
//		),
//	)
//	s.Require().NoError(err)
//	s.Require().NoError(fixture.Load())
//
//	s.Token = s.TS.Login("admin@example.com", "securepass")
//}
//
//func (s *GenerateTokensTestSuite) TearDownTest() {
//	s.TS.TruncateAll()
//}
//
//func (s *GenerateTokensTestSuite) TestGenerateTokensSuccess() {
//	const (
//		testPlaceID = "a8c52b0c-8f11-4b9c-9c3f-123456789abc"
//		count       = 3
//	)
//
//	payload := map[string]any{
//		"place_id": testPlaceID,
//		"count":    count,
//	}
//	body, _ := json.Marshal(payload)
//
//	req := httptest.NewRequest(http.MethodPost, "/admin/tokens", bytes.NewReader(body))
//	req.Header.Set("Authorization", "Bearer "+s.Token)
//	req.Header.Set("Content-Type", "application/json")
//
//	rec := httptest.NewRecorder()
//	s.TS.App.ServeHTTP(rec, req)
//
//	s.Require().Equal(http.StatusOK, rec.Code)
//
//	var resp struct {
//		Tokens []string `json:"tokens"`
//	}
//	err := json.NewDecoder(rec.Body).Decode(&resp)
//	s.Require().NoError(err)
//	s.Require().Equal(count, len(resp.Tokens))
//}
//
//func (s *GenerateTokensTestSuite) TestGenerateTokensForbiddenForUser() {
//	const (
//		testPlaceID = "a8c52b0c-8f11-4b9c-9c3f-123456789abc"
//		count       = 2
//	)
//
//	s.Token = s.TS.Login("bob@example.com", "password123")
//
//	payload := map[string]any{
//		"place_id": testPlaceID,
//		"count":    count,
//	}
//	body, _ := json.Marshal(payload)
//
//	req := httptest.NewRequest(http.MethodPost, "/admin/tokens", bytes.NewReader(body))
//	req.Header.Set("Authorization", "Bearer "+s.Token)
//	req.Header.Set("Content-Type", "application/json")
//
//	rec := httptest.NewRecorder()
//	s.TS.App.ServeHTTP(rec, req)
//
//	s.Require().Equal(http.StatusForbidden, rec.Code)
//	s.Require().Contains(rec.Body.String(), "only admin can generate tokens")
//}
//
//func (s *GenerateTokensTestSuite) TestGenerateTokensInvalidInput() {
//	const (
//		testPlaceID = "a8c52b0c-8f11-4b9c-9c3f-123456789abc"
//		count       = 0
//	)
//
//	payload := map[string]any{
//		"place_id": testPlaceID,
//		"count":    count,
//	}
//	body, _ := json.Marshal(payload)
//
//	req := httptest.NewRequest(http.MethodPost, "/admin/tokens", bytes.NewReader(body))
//	req.Header.Set("Authorization", "Bearer "+s.Token)
//	req.Header.Set("Content-Type", "application/json")
//
//	rec := httptest.NewRecorder()
//	s.TS.App.ServeHTTP(rec, req)
//
//	s.Require().Equal(http.StatusBadRequest, rec.Code)
//	s.Require().Contains(rec.Body.String(), "invalid input")
//}
