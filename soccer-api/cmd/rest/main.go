package main

import (
	"fmt"

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

	mailService := service.NewEmailService(cfg.MailNotification)
	sendMail := false
	if sendMail {
		err := mailService.SendEmail(
			"doug_21@gol.com",
			"Seja Bem Vindo ao SoccerFan",
			"Estamos gratos por vc ter ativado sua conta em nosso site! Aproveite aos jogos do seu time favorito.",
		)
		if err != nil {
			fmt.Print(err)
		}
	}

	err := routes.NewGinRoutes(competitionService, fanService).Run()
	if err != nil {
		panic("cant initiate routes")
	}
}
