package routes

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
)

type ginRoutesSuite struct {
	suite.Suite
	r *ginRoutes
}

func TestGinRoutesSuite(t *testing.T) {
	suite.Run(t, new(ginRoutesSuite))
}

func (suite *ginRoutesSuite) SetupSuite() {
	fDB := NewFakeDB()
	fakeCRepo := NewFakeChampionshipRepo(fDB)
	fakeFRepo := NewFakeFanRepo(fDB)
	cService := service.NewChampionship(fakeCRepo)
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
