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
	s.TS.TruncateAll()

	db := stdlib.OpenDBFromPool(s.TS.DB)
	defer db.Close()

	fixture, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Files(
			"../fixtures/users.yml",
			"../fixtures/places.yml",
		),
	)
	s.Require().NoError(err, "init fixtures failed")
	s.Require().NoError(fixture.Load(), "load fixtures failed")

	s.Token = s.TS.Login("admin@example.com", "securepass")
}

func (s *CreatePlaceTestSuite) TearDownTest() {
	s.TS.TruncateAll()
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

	s.Require().Equal(http.StatusCreated, rec.Code)
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

	s.Require().Equal(http.StatusBadRequest, rec.Code)
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

	s.Require().Equal(http.StatusUnauthorized, rec.Code)
}
