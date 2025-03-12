package routes

import (
	"bytes"
	"context"
	"fmt"
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
	competitions []*entity.Competition
	teams        []*entity.Team
	matches      []*entity.Match
	fans         []*entity.Fan
}

func NewFakeDB() *fakeDB {

	scoreOne := 1
	scoreTwo := 2

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
		competitions: []*entity.Competition{
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

		matches: []*entity.Match{
			{
				ID:            3,
				UID:           uuid.MustParse("00000000-0000-0000-0000-000000003003"),
				Round:         2,
				CompetitionID: 1,
				HomeTeam:      flamengo,
				AwayTeam:      corinthians,
				HomeTeamScore: nil,
				AwayTeamScore: nil,
			},
			{
				ID:            2,
				UID:           uuid.MustParse("00000000-0000-0000-0000-000000003002"),
				Round:         1,
				CompetitionID: 1,
				HomeTeam:      santos,
				AwayTeam:      corinthians,
				HomeTeamScore: &scoreTwo,
				AwayTeamScore: &scoreTwo,
			},
			{
				ID:            1,
				UID:           uuid.MustParse("00000000-0000-0000-0000-000000003001"),
				Round:         1,
				CompetitionID: 1,
				HomeTeam:      flamengo,
				AwayTeam:      vasco,
				HomeTeamScore: &scoreTwo,
				AwayTeamScore: &scoreOne,
			},
		},
	}
}

type fakeCompetitionRepo struct {
	dbConn *fakeDB
}

func NewFakeCompetitionRepo(dbConn *fakeDB) *fakeCompetitionRepo {
	return &fakeCompetitionRepo{
		dbConn,
	}
}

func (fcr *fakeCompetitionRepo) FindAll(_ context.Context) ([]*entity.Competition, error) {
	return fcr.dbConn.competitions, nil
}

func (fcr *fakeCompetitionRepo) FindMatchsByCompetitionUID(_ context.Context, uid uuid.UUID) ([]*entity.Match, error) {
	return fcr.dbConn.matches, nil
}

type fakeFanRepo struct {
	dbConn *fakeDB
}

func NewFakeFanRepo(dbConn *fakeDB) *fakeFanRepo {
	return &fakeFanRepo{
		dbConn,
	}
}

func (ffr *fakeFanRepo) Create(_ context.Context, fEntity *entity.Fan) (*entity.Fan, error) {

	var teamFounded *entity.Team
	for _, team := range ffr.dbConn.teams {
		if team.Name == fEntity.Team.Name {
			teamFounded = team
		}
	}

	if teamFounded == nil {
		return nil, fmt.Errorf("team not found")
	}

	fanCreated := &entity.Fan{
		ID:    len(ffr.dbConn.fans) + 1,
		UID:   fEntity.UID,
		Name:  fEntity.Name,
		Email: fEntity.Email,
		Team:  teamFounded,
	}

	ffr.dbConn.fans = append(ffr.dbConn.fans, fanCreated)
	return fanCreated, nil
}

type ginRoutesSuite struct {
	suite.Suite
	r *ginRoutes
}

func TestGinRoutesSuite(t *testing.T) {
	suite.Run(t, new(ginRoutesSuite))
}

func (suite *ginRoutesSuite) SetupSuite() {
	fDB := NewFakeDB()
	fakeCRepo := NewFakeCompetitionRepo(fDB)
	fakeFRepo := NewFakeFanRepo(fDB)
	cService := service.NewCompetition(fakeCRepo)
	fService := service.NewFan(fakeFRepo)
	suite.r = NewGinRoutes(cService, fService)
}

func (suite *ginRoutesSuite) TestGetChampionshipsSuccesfully() {
	req, err := http.NewRequest("GET", "/campeonatos", nil)
	assert.NoError(suite.T(), err)

	resp := httptest.NewRecorder()
	suite.r.engine.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	championship := gjson.Get(resp.Body.String(), "campeonatos").Array()
	assert.Equal(suite.T(), 2, len(championship))
}

func (suite *ginRoutesSuite) TestGetChampionshipMatchsWithoutFiltersSuccesfully() {
	req, err := http.NewRequest("GET", "/campeonatos/00000000-0000-0000-0000-000000001001/partidas", nil)
	assert.NoError(suite.T(), err)

	resp := httptest.NewRecorder()
	suite.r.engine.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	rounds := gjson.Get(resp.Body.String(), "rodadas").Array()
	assert.Equal(suite.T(), 2, len(rounds))
}

func (suite *ginRoutesSuite) TestPostFanSuccessfully() {
	reqBody := `
		{
			"nome" : "João Silva",
			"email": "joao.silva@example.com",
			"time" : "Flamengo"
		}
	`

	req, err := http.NewRequest("POST", "/torcedores", bytes.NewBuffer([]byte(reqBody)))
	assert.NoError(suite.T(), err)

	resp := httptest.NewRecorder()
	suite.r.engine.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusAccepted, resp.Code)

	assert.Equal(suite.T(), "João Silva", gjson.Get(resp.Body.String(), "nome").String())
	assert.Equal(suite.T(), "joao.silva@example.com", gjson.Get(resp.Body.String(), "email").String())
	assert.Equal(suite.T(), "Flamengo", gjson.Get(resp.Body.String(), "time").String())
	assert.Equal(
		suite.T(),
		"Cadastro realizado com sucesso",
		gjson.Get(resp.Body.String(), "mensagem").String(),
	)
}
