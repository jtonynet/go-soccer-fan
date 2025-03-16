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
	competitionRepo := gormrepo.NewCompetition(gormConn)
	fanRepo := gormrepo.NewFan(gormConn)

	competitionService := service.NewCompetition(competitionRepo)
	fanService := service.NewFan(fanRepo)

	teamRepo := gormrepo.NewTeam(gormConn)
	broadcastService := service.NewEmailBroadcast(
		cfg.MailNotification,
		teamRepo,
	)

	err := routes.NewGinRoutes(
		competitionService,
		fanService,
		broadcastService,
	).Run()
	if err != nil {
		panic("can't start routes")
	}
}
