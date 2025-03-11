package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginRouter struct {
	Router *gin.Engine
}

func newGinRouter() *ginRouter {
	e := gin.Default()

	e.GET("/campeonatos", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	return &ginRouter{
		Router: e,
	}
}
