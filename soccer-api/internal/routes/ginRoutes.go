package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/service"
)

// TODO: Código fake de testes apenas para validar a API REST Dockerizada. Será removido no próximo PR
type fakeDB struct {
	championships []*entity.Championship
	teams         []*entity.Team
	matchs        []*entity.Match
	fans          []*entity.Fan
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

func NewFakeChampionshipRepo(dbConn *fakeDB) *fakeChampionshipRepo {
	return &fakeChampionshipRepo{
		dbConn,
	}
}

func (fcr *fakeChampionshipRepo) FindAll(_ context.Context) ([]*entity.Championship, error) {
	return fcr.dbConn.championships, nil
}

func (fcr *fakeChampionshipRepo) FindMatchsByChampionshipUID(_ context.Context, uid uuid.UUID) ([]*entity.Match, error) {
	return fcr.dbConn.matchs, nil
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

type ginRoutes struct {
	engine *gin.Engine
}

func NewGinRoutes(cService *service.Championship, fService *service.Fan) *ginRoutes {
	e := gin.Default()

	e.GET("/campeonatos", func(c *gin.Context) {
		result, _ := cService.FindAll()
		c.JSON(http.StatusOK, result)
	})

	e.GET("/campeonatos/:uid/partidas", func(c *gin.Context) {
		uid, err := uuid.Parse(c.Param("uid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": "campeonato com URI ID inválido",
			})
			return
		}

		result, _ := cService.FindMatchsByChampionshipUID(uid)
		c.JSON(http.StatusOK, result)
	})

	e.POST("/torcedores", func(c *gin.Context) {
		var fReq dto.FanCreateRequest
		if err := c.ShouldBindJSON(&fReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": "corpo da requisição inválido",
			})
			return
		}

		// TODO: VALIDATES DTO IN FUTURE

		fResp, err := fService.Create(&fReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"erro": "erro interno, tente novamente mais tarde",
			})
			return
		}

		c.JSON(http.StatusAccepted, fResp)
	})

	return &ginRoutes{
		engine: e,
	}
}

func (gr *ginRoutes) Run() error {
	return gr.engine.Run()
}
