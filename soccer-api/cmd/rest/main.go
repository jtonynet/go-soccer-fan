package main

import (
	"log"

	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/database"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/pubsub"
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

	pubSub, err := pubsub.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("can't instantiate pubsub: %v", err)
	}

	teamRepo := gormrepo.NewTeam(gormConn)
	broadcastService := service.NewBroadcast(pubSub, teamRepo, cfg.RabbitMQ.MatchNoticationsQueue)

	err = routes.NewGinRoutes(
		competitionService,
		fanService,
		broadcastService,
	).Run()
	if err != nil {
		panic("can't start routes")
	}
}
