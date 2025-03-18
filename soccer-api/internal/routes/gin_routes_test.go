package routes

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/database"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/pubsub"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository/gormrepo"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

var bearerToken string

type ginRoutesSuite struct {
	suite.Suite
	r  *ginRoutes
	db *gorm.DB
}

func TestGinRoutesSuite(t *testing.T) {
	suite.Run(t, new(ginRoutesSuite))
}

func (suite *ginRoutesSuite) setupDB(cfg *config.Database) {
	gormConn, err := database.NewGormCom(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db := gormConn.GetDB()

	suite.db = db

	tx := db.Begin()

	if err := tx.Exec(`
		INSERT INTO public.users (created_at, updated_at, deleted_at, uid, username, "password", "name", email) VALUES
			(NOW(), NOW(), NULL, '00000000-0000-0000-0000-000000000001', 'admin', '$2a$10$MrnQU3LS5qtk9Ca.EETK1.Yj4M4AQoNQAQ08gEu9yl4q4lr9lW9gO', 'Edson Arantes do Nascimento', 'pele@soccerfan.com');
	`).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert into users: %v", err)
	}

	if err := tx.Exec(`
		INSERT INTO public.competitions (created_at, updated_at, deleted_at, uid, external_id, "name", season) VALUES
			(NOW(), NOW(), NULL, '00000000-0000-0000-0000-000000001001', '2013', 'Campeonato Brasileiro', '2025'),
			(NOW(), NOW(), NULL, '00000000-0000-0000-0000-000000001002', '2025', 'UEFA Champions League', '2025');
	`).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert into competitions: %v", err)
	}

	if err := tx.Exec(`
		INSERT INTO public.teams (created_at, updated_at, deleted_at, uid, external_id, "name", full_name) VALUES
			(NOW(), NOW(), NULL, '00000000-0000-0000-0000-000000002001', '61', 'Flamengo', 'Flamengo FC'),
			(NOW(), NOW(), NULL, '00000000-0000-0000-0000-000000002002', '62', 'Vasco', 'Vasco FC'),
			(NOW(), NOW(), NULL, '00000000-0000-0000-0000-000000002003', '63', 'Santos', 'Santos FC'),
			(NOW(), NOW(), NULL, '00000000-0000-0000-0000-000000002004', '64', 'Corinthians', 'Corinthians FC');
	`).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert into teams: %v", err)
	}

	if err := tx.Exec(`
		INSERT INTO public.matches (created_at, updated_at, deleted_at, uid, external_id, competition_id, home_team_id, away_team_id, round, home_team_score, away_team_score) VALUES
			(NOW(), NOW(), NULL, '00000000-0000-0000-0000-000000003001', '428747', 1, 1, 2, 1, NULL, NULL),
			(NOW(), NOW(), NULL, '00000000-0000-0000-0000-000000003002', '428748', 1, 3, 4, 1, NULL, NULL),
			(NOW(), NOW(), NULL, '00000000-0000-0000-0000-000000003003', '428749', 1, 1, 3, 2, NULL, NULL);
	`).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert into matches: %v", err)
	}

	tx.Commit()
}

func (suite *ginRoutesSuite) TearDownSuite() {
	tx := suite.db.Begin()
	if err := suite.db.Exec(`
		TRUNCATE TABLE public.users  RESTART IDENTITY CASCADE;
		TRUNCATE TABLE public.competitions RESTART IDENTITY CASCADE;
		TRUNCATE TABLE public.teams RESTART IDENTITY CASCADE;
		TRUNCATE TABLE public.fans RESTART IDENTITY CASCADE;
		TRUNCATE TABLE public.matches RESTART IDENTITY CASCADE;
		
		ALTER SEQUENCE public.users_id_seq RESTART WITH 1;
		ALTER SEQUENCE public.competitions_id_seq RESTART WITH 1;
		ALTER SEQUENCE public.fans_id_seq RESTART WITH 1;
		ALTER SEQUENCE public.matches_id_seq RESTART WITH 1;
		ALTER SEQUENCE public.teams_id_seq RESTART WITH 1;
	`).Error; err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert clean test database: %v", err)
	}
	tx.Commit()
}

func (suite *ginRoutesSuite) SetupSuite() {
	cfg := config.LoadConfig()
	cfg.Database.DBname = "soccer_db_test"

	suite.setupDB(cfg.Database)

	pubSub, err := pubsub.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("can't connect to pubsub: %v", err)
	}

	gormConn, err := database.NewGormCom(cfg.Database)
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}

	userRepo := gormrepo.NewUser(gormConn)
	competitionRepo := gormrepo.NewCompetition(gormConn)
	fanRepo := gormrepo.NewFan(gormConn)
	teamRepo := gormrepo.NewTeam(gormConn)

	userService := service.NewUser(userRepo)
	competitionService := service.NewCompetition(competitionRepo)
	fanService := service.NewFan(fanRepo)
	broadcastService := service.NewBroadcast(
		pubSub,
		teamRepo,
		cfg.RabbitMQ.MatchNoticationsQueue,
	)

	suite.r = NewGinRoutes(
		userService,
		competitionService,
		fanService,
		broadcastService,
	)

	reqBody := `
		{
		  "usuario": "admin",
		  "senha": "admin"
		}
	`

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte(reqBody)))

	resp := httptest.NewRecorder()
	suite.r.engine.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusAccepted, resp.Code)

	bearerToken = gjson.Get(resp.Body.String(), "token").String()
}

func (suite *ginRoutesSuite) TestGetChampionshipsSuccesfully() {
	req, err := http.NewRequest("GET", "/campeonatos", nil)
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp := httptest.NewRecorder()
	suite.r.engine.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	championship := gjson.Get(resp.Body.String(), "campeonatos").Array()
	assert.Equal(suite.T(), 2, len(championship))
}

func (suite *ginRoutesSuite) TestGetChampionshipsUnauthorized() {
	req, err := http.NewRequest("GET", "/campeonatos", nil)
	assert.NoError(suite.T(), err)

	resp := httptest.NewRecorder()
	suite.r.engine.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
}

func (suite *ginRoutesSuite) TestGetChampionshipMatchsWithoutFiltersSuccesfully() {
	req, err := http.NewRequest("GET", "/campeonatos/00000000-0000-0000-0000-000000001001/partidas", nil)
	assert.NoError(suite.T(), err)

	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp := httptest.NewRecorder()
	suite.r.engine.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusOK, resp.Code)

	rounds := gjson.Get(resp.Body.String(), "rodadas").Array()
	assert.Equal(suite.T(), 2, len(rounds))
}

func (suite *ginRoutesSuite) TestGetChampionshipMatchsWithoutFiltersUnauthorized() {
	req, err := http.NewRequest("GET", "/campeonatos/00000000-0000-0000-0000-000000001001/partidas", nil)
	assert.NoError(suite.T(), err)

	resp := httptest.NewRecorder()
	suite.r.engine.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
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

func (suite *ginRoutesSuite) TestPostFanNameTooShort() {
	reqBody := `
		{
			"nome" : "Jo",
			"email": "joao.silva@example.com",
			"time" : "Flamengo"
		}
	`

	req, err := http.NewRequest("POST", "/torcedores", bytes.NewBuffer([]byte(reqBody)))
	assert.NoError(suite.T(), err)

	resp := httptest.NewRecorder()
	suite.r.engine.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *ginRoutesSuite) TestPostFanInvalidEmail() {
	reqBody := `
		{
			"nome" : "João Silva",
			"email": "invalid-email",
			"time" : "Flamengo"
		}
	`

	req, err := http.NewRequest("POST", "/torcedores", bytes.NewBuffer([]byte(reqBody)))
	assert.NoError(suite.T(), err)

	resp := httptest.NewRecorder()
	suite.r.engine.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *ginRoutesSuite) TestPostFanTeamNameTooShort() {
	reqBody := `
		{
			"nome" : "João Silva",
			"email": "joao.silva@example.com",
			"time" : "Flam"
		}
	`

	req, err := http.NewRequest("POST", "/torcedores", bytes.NewBuffer([]byte(reqBody)))
	assert.NoError(suite.T(), err)

	resp := httptest.NewRecorder()
	suite.r.engine.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}

func (suite *ginRoutesSuite) TestPostFanMissingFields() {
	reqBody := `
		{
			"nome" : "",
			"email": "",
			"time" : ""
		}
	`

	req, err := http.NewRequest("POST", "/torcedores", bytes.NewBuffer([]byte(reqBody)))
	assert.NoError(suite.T(), err)

	resp := httptest.NewRecorder()
	suite.r.engine.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
}
