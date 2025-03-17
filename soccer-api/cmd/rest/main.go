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

	pubSub, err := pubsub.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("can't connect to pubsub: %v", err)
	}

	gormConn, err := database.NewGormCom(cfg.Database)
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}

	userRepo := gormrepo.NewUser(gormConn)
	competitionRepo := gormrepo.NewCompetition(gormConn)
	fanRepo := gormrepo.NewFan(gormConn)
	teamRepo := gormrepo.NewTeam(gormConn)

	userService := service.NewUser(userRepo)
	competitionService := service.NewCompetition(competitionRepo)
	fanService := service.NewFan(fanRepo)
	broadcastService := service.NewBroadcast(
		pubSub,
		teamRepo,
		cfg.RabbitMQ.MatchNoticationsQueue,
	)

	err = routes.NewGinRoutes(
		userService,
		competitionService,
		fanService,
		broadcastService,
	).Run()
	if err != nil {
		panic("can't start routes")
	}
}
