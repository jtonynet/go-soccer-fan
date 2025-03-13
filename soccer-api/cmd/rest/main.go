package main

import (
	"context"

	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/database"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository/gormrepo"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/routes"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	gormConn := database.NewGormCom(cfg.Database)
	competitionRepo := gormrepo.NewCompetition(gormConn)
	fanRepo := gormrepo.NewFan(gormConn)

	competitionService := service.NewCompetition(competitionRepo)
	fanService := service.NewFan(fanRepo)

	// TODO: mover para cmd/scheduler/main
	areaRepo := gormrepo.NewArea(gormConn)
	fetchService := service.NewDataFetchService(
		cfg.ExternalApi,
		areaRepo,
		competitionRepo,
	)
	fetchService.FetchAndStoreAreaData(context.Background())

	err := routes.NewGinRoutes(competitionService, fanService).Run()
	if err != nil {
		panic("cant initiate routes")
	}
}
