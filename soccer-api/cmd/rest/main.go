package main

import (
	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/database"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository/gormrepo"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/routes"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	gormConn := database.NewGormCom(cfg.Database)
	championshipRepo := gormrepo.NewChampionship(gormConn)
	fanRepo := gormrepo.NewFan(gormConn)

	championshipService := service.NewChampionship(championshipRepo)
	fanService := service.NewFan(fanRepo)

	err := routes.NewGinRoutes(championshipService, fanService).Run()
	if err != nil {
		panic("cant initiate routes")
	}
}
