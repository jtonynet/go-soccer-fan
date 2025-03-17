package main

import (
	"context"
	"log"
	"os"

	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/database"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository/gormrepo"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	gormConn, err := database.NewGormCom(cfg.Database)
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}

	competitionRepo := gormrepo.NewCompetition(gormConn)
	teamRepo := gormrepo.NewTeam(gormConn)
	matchRepo := gormrepo.NewMatch(gormConn)
	dataFetchService := service.NewDataFetchService(
		cfg.ExternalApi,
		competitionRepo,
		teamRepo,
		matchRepo,
	)

	if len(os.Args) < 2 {
		log.Fatalf("Comando esperado: `import competitions`")
	}

	command := os.Args[1]

	switch command {
	case "import":
		if len(os.Args) < 3 {
			log.Fatalf("Comando esperado: `import competitions`")
		}
		subCommand := os.Args[2]
		switch subCommand {
		case "competitions":
			if err := dataFetchService.FetchAndStore(context.Background()); err != nil {
				log.Fatalf("Erro ao buscar e armazenar dados: %v", err)
			}
			log.Println("Importação de competições concluída com sucesso!")
		default:
			log.Fatalf("Comando desconhecido: %s", subCommand)
		}
	default:
		log.Fatalf("Comando desconhecido: %s", command)
	}
}
