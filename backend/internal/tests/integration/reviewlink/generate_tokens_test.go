package reviewlink

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/kulikovroman08/reviewlink-backend/internal/tests/integration"

	"github.com/stretchr/testify/suite"
)

type GenerateTokensTestSuite struct {
	suite.Suite
	TS    *integration.TestSetup
	Token string
}

func TestGenerateTokensSuite(t *testing.T) {
	suite.Run(t, new(GenerateTokensTestSuite))
}

func (s *GenerateTokensTestSuite) SetupSuite() {
	s.TS = integration.NewTestSetup()
}

func (s *GenerateTokensTestSuite) TearDownSuite() {
	s.TS.Close()
}

func (s *GenerateTokensTestSuite) SetupTest() {
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
			"../fixtures/places.yml",
		),
	)
	require.NoError(s.T(), err)
	require.NoError(s.T(), fixture.Load())

	s.Token = s.TS.Login("admin@example.com", "securepass")
}

func (s *GenerateTokensTestSuite) TestGenerateTokensSuccess() {
	const (
		testPlaceID = "a8c52b0c-8f11-4b9c-9c3f-123456789abc"
		count       = 3
	)

	payload := map[string]any{
		"place_id": testPlaceID,
		"count":    count,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/admin/tokens", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	s.Require().Equal(http.StatusOK, rec.Code)

	var resp struct {
		Tokens []string `json:"tokens"`
	}
	err := json.NewDecoder(rec.Body).Decode(&resp)

	require.NoError(s.T(), err)
	require.Equal(s.T(), count, len(resp.Tokens))
}

func (s *GenerateTokensTestSuite) TestGenerateTokensForbiddenForUser() {
	const (
		testPlaceID  = "a8c52b0c-8f11-4b9c-9c3f-123456789abc"
		count        = 2
		errOnlyAdmin = `{"error":"only admin can generate tokens"}`
	)

	s.Token = s.TS.Login("bob@example.com", "password123")

	payload := map[string]any{
		"place_id": testPlaceID,
		"count":    count,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/admin/tokens", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusForbidden, rec.Code)
	require.Equal(s.T(), errOnlyAdmin, rec.Body.String())
}

func (s *GenerateTokensTestSuite) TestGenerateTokensInvalidInput() {
	const (
		testPlaceID   = "a8c52b0c-8f11-4b9c-9c3f-123456789abc"
		count         = 0
		errInvalidInp = `{"error":"invalid input"}`
	)

	payload := map[string]any{
		"place_id": testPlaceID,
		"count":    count,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/admin/tokens", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusBadRequest, rec.Code)
	require.Equal(s.T(), errInvalidInp, rec.Body.String())
}
