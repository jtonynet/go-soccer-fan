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
	teams         []*entity.Team
	matchs        []*entity.Match
}

func NewFakeDB() *fakeDB {

	scoreOne := 1
	scoreTwo := 3

	flamengo := &entity.Team{
		ID:   1,
		UID:  uuid.MustParse("00000000-0000-0000-0000-000000002001"),
		Name: "Flamengo",
	}

	vasco := &entity.Team{
		ID:   2,
		UID:  uuid.MustParse("00000000-0000-0000-0000-000000002002"),
		Name: "Vasco",
	}

	santos := &entity.Team{
		ID:   3,
		UID:  uuid.MustParse("00000000-0000-0000-0000-000000002003"),
		Name: "Santos",
	}

	corinthians := &entity.Team{
		ID:   4,
		UID:  uuid.MustParse("00000000-0000-0000-0000-000000002004"),
		Name: "Corinthians",
	}

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

		teams: []*entity.Team{
			flamengo,
			vasco,
			santos,
			corinthians,
		},

		matchs: []*entity.Match{
			{
				ID:             3,
				UID:            uuid.MustParse("00000000-0000-0000-0000-000000003003"),
				Round:          2,
				ChampionshipID: 1,
				HomeTeam:       flamengo,
				AwayTeam:       corinthians,
				HomeTeamScore:  nil,
				AwayTeamScore:  nil,
			},
			{
				ID:             2,
				UID:            uuid.MustParse("00000000-0000-0000-0000-000000003002"),
				Round:          1,
				ChampionshipID: 1,
				HomeTeam:       santos,
				AwayTeam:       corinthians,
				HomeTeamScore:  &scoreTwo,
				AwayTeamScore:  &scoreTwo,
			},
			{
				ID:             1,
				UID:            uuid.MustParse("00000000-0000-0000-0000-000000003001"),
				Round:          1,
				ChampionshipID: 1,
				HomeTeam:       flamengo,
				AwayTeam:       vasco,
				HomeTeamScore:  &scoreTwo,
				AwayTeamScore:  &scoreOne,
			},
		},
	}
}

type fakeChampionshipRepo struct {
	dbConn *fakeDB
}

func NewfakeChampionshipRepo(dbConn *fakeDB) *fakeChampionshipRepo {
	return &fakeChampionshipRepo{
		dbConn,
	}
}

func (fcr *fakeChampionshipRepo) FindAll(ctx context.Context) ([]*entity.Championship, error) {
	return fcr.dbConn.championships, nil
}

func (fcr *fakeChampionshipRepo) FindMatchsByChampionshipUID(ctx context.Context, uid uuid.UUID) ([]*entity.Match, error) {
	return fcr.dbConn.matchs, nil
}

type ginRoutesSuite struct {
	suite.Suite
	r *ginRouter
}

func TestGinRoutesSuite(t *testing.T) {
	suite.Run(t, new(ginRoutesSuite))
}

func (suite *ginRoutesSuite) SetupSuite() {
	fDB := NewFakeDB()
	fakeCRepo := NewfakeChampionshipRepo(fDB)
	cService := service.NewChampionship(fakeCRepo)
	suite.r = NewGinRouter(cService)
}

func (suite *ginRoutesSuite) TestGetChampionshipsSuccesfully() {
	req, err := http.NewRequest("GET", "/campeonatos", nil)
	assert.NoError(suite.T(), err)

	resp := httptest.NewRecorder()
	suite.r.Router.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	c := gjson.Get(resp.Body.String(), "campeonatos").Array()
	assert.Equal(suite.T(), 2, len(c))
}

func (suite *ginRoutesSuite) TestGetChampionshipMatchsWithoutFiltersSuccesfully() {
	req, err := http.NewRequest("GET", "/campeonatos/00000000-0000-0000-0000-000000001001/partidas", nil)
	assert.NoError(suite.T(), err)

	resp := httptest.NewRecorder()
	suite.r.Router.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	c := gjson.Get(resp.Body.String(), "rodadas").Array()
	assert.Equal(suite.T(), 2, len(c))
}
