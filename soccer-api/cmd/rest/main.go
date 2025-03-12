package main

import (
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/routes"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/service"
)

func main() {

	// TODO: Código chamando fake de testes apenas para validar a API REST Dockerizada. Será removido no próximo PR
	dbConn := routes.NewFakeDB()

	cRepo := routes.NewFakeChampionshipRepo(dbConn)
	cService := service.NewChampionship(cRepo)

	fRepo := routes.NewFakeFanRepo(dbConn)
	fService := service.NewFan(fRepo)

	err := routes.NewGinRoutes(cService, fService).Run()
	if err != nil {
		panic("cant initiate routes")
	}
}
