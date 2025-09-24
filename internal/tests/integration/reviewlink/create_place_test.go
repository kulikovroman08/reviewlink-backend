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

type CreatePlaceTestSuite struct {
	suite.Suite
	TS    *integration.TestSetup
	Token string
}

func TestCreatePlaceSuite(t *testing.T) {
	suite.Run(t, new(CreatePlaceTestSuite))
}

func (s *CreatePlaceTestSuite) SetupSuite() {
	s.TS = integration.NewTestSetup()
}

func (s *CreatePlaceTestSuite) TearDownSuite() {
	s.TS.Close()
}

func (s *CreatePlaceTestSuite) SetupTest() {
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
	require.NoError(s.T(), err, "init fixtures failed")
	require.NoError(s.T(), fixture.Load(), "load fixtures failed")

	s.Token = s.TS.Login("admin@example.com", "securepass")
}

func (s *CreatePlaceTestSuite) TestSuccessCreatePlace() {
	payload := map[string]string{
		"name":    "My Cafe",
		"address": "123 Main St",
	}
	data, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/places", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.Token)

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusCreated, rec.Code)
}

func (s *CreatePlaceTestSuite) TestCreatePlaceMissingName() {
	payload := map[string]string{
		"address": "123 Main St",
	}
	data, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/places", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.Token)

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusBadRequest, rec.Code)
}

func (s *CreatePlaceTestSuite) TestCreatePlaceUnauthorized() {
	payload := map[string]string{
		"name":    "Cafe",
		"address": "123 Main St",
	}
	data, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/places", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	s.TS.App.ServeHTTP(rec, req)

	require.Equal(s.T(), http.StatusUnauthorized, rec.Code)
}
