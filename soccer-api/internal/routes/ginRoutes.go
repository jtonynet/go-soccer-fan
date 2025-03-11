package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	return &ginRouter{
		Router: e,
	}
}

func (gr *ginRouter) Run() error {
	return gr.Router.Run()
}
