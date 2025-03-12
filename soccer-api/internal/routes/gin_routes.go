package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/service"
)

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
