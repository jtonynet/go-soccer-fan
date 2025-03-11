package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/service"
)

type ginRouter struct {
	Router *gin.Engine
}

func NewGinRouter(cService *service.Championship) *ginRouter {
	e := gin.Default()

	e.GET("/campeonatos", func(c *gin.Context) {
		result, _ := cService.FindAll()
		c.JSON(http.StatusOK, result)
	})

	e.GET("/campeonatos/:uid/partidas", func(c *gin.Context) {
		uid, err := uuid.Parse(c.Param("uid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "campeonato com ID invalido",
			})
			return
		}

		result, _ := cService.FindMatchsByChampionshipUID(uid)
		c.JSON(http.StatusOK, result)
	})

	return &ginRouter{
		Router: e,
	}
}

func (gr *ginRouter) Run() error {
	return gr.Router.Run()
}
