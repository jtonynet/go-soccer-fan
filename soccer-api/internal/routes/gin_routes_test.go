package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
)

type fakeDB struct {
	championships []*entity.Championship
}

func getFakeDB() *fakeDB {
	return &fakeDB{
		championships: []*entity.Championship{
			{
				ID:     1,
				UID:    uuid.MustParse("00000000-0000-0000-0000-000000001001"),
				Name:   "Campeonato Brasileiro",
				Season: "2025",
			},
			{
				ID:     2,
				UID:    uuid.MustParse("00000000-0000-0000-0000-000000001002"),
				Name:   "UEFA Champions League",
				Season: "2025",
			},
		},
	}
}

type fakeChampionshipRepo struct {
	dbConn *fakeDB
}

func newfakeChampionshipRepo(dbConn *fakeDB) *fakeChampionshipRepo {
	return &fakeChampionshipRepo{
		dbConn,
	}
}

func (fcr *fakeChampionshipRepo) FindAll(ctx context.Context) ([]*entity.Championship, error) {
	return nil, nil
}

type ginRoutesSuite struct {
	suite.Suite
	r *ginRouter
}

func TestGinRoutesSuite(t *testing.T) {
	suite.Run(t, new(ginRoutesSuite))
}

func (suite *ginRoutesSuite) SetupSuite() {
	fDB := getFakeDB()
	fakeCRepo := newfakeChampionshipRepo(fDB)
	cService := service.NewChampionship(fakeCRepo)
	suite.r = NewGinRouter(cService)
}

func (suite *ginRoutesSuite) TestHappyPath() {
	req, err := http.NewRequest("GET", "/campeonatos", nil)
	assert.NoError(suite.T(), err)

	resp := httptest.NewRecorder()
	suite.r.Router.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	c := gjson.Get(resp.Body.String(), "campeonatos").Array()
	assert.Equal(suite.T(), 2, len(c))
}
