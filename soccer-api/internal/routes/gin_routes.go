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

func NewGinRoutes(
	competitionService *service.Competition,
	fanService *service.Fan,
	broadcastService *service.Broadcast,
) *ginRoutes {
	e := gin.Default()

	e.GET("/campeonatos", func(c *gin.Context) {
		result, _ := competitionService.FindAll()
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

		result, err := competitionService.FindMatchsByCompetitionUID(uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"erro": "erro interno, tente novamente mais tarde",
			})
		}
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

		fResp, err := fanService.Create(&fReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"erro": "erro interno, tente novamente mais tarde",
			})
			return
		}

		c.JSON(http.StatusAccepted, fResp)
	})

	e.POST("/broadcast", func(c *gin.Context) {
		var bReq dto.BroadcastSendRequest
		if err := c.ShouldBindJSON(&bReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"erro": "corpo da requisição inválido",
			})
			return
		}

		// TODO: VALIDATES DTO IN FUTURE

		bResp, err := broadcastService.Send(&bReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"erro": "erro interno, tente novamente mais tarde",
			})
			return
		}

		c.JSON(http.StatusAccepted, bResp)

	})

	return &ginRoutes{
		engine: e,
	}
}

func (gr *ginRoutes) Run() error {
	return gr.engine.Run()
}
