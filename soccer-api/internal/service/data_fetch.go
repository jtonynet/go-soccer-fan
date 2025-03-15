package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository"
)

var retryInitialInterval = time.Duration(1)
var retryMaxElapsedTime = time.Duration(60)

var workerPoolSize = 5

var teamsProcessedCount = 0
var matchesProcessedCount = 0

var competitionsCached sync.Map
var teamsCached sync.Map

type DataFetch struct {
	cfg             *config.ExternalApi
	competitionRepo repository.Competition
	teamRepo        repository.Team
	matchRepo       repository.Match
	httpClient      *http.Client
	startTime       time.Time
}

func NewDataFetchService(
	cfg *config.ExternalApi,
	competitionRepo repository.Competition,
	teamRepo repository.Team,
	matchRepo repository.Match,
) *DataFetch {
	return &DataFetch{
		cfg:             cfg,
		competitionRepo: competitionRepo,
		teamRepo:        teamRepo,
		matchRepo:       matchRepo,
		httpClient:      &http.Client{Timeout: 10 * time.Second},
		startTime:       time.Now(),
	}
}

func (s *DataFetch) FetchAndStore(ctx context.Context) error {
	uriFragment := "v4/competitions"
	body, err := s.fetchFromExternalAPI(ctx, uriFragment)
	if err != nil {
		return err
	}

	var response dto.CompetitionResponseListExternalAPI
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	var competitionEntities []*entity.Competition
	for _, cResponse := range response.Competitions {
		competitionEntities = append(competitionEntities, &entity.Competition{
			UID:        uuid.New(),
			ExternalId: fmt.Sprintf("%d", cResponse.ID),
			Name:       cResponse.Name,
			Season:     fmt.Sprintf("%v", time.Now().Year()),
		})
	}

	competitions, err := s.competitionRepo.CreateOrUpdateInBatch(ctx, competitionEntities)
	if err != nil {
		log.Printf("Erro ao salvar competições: %v", err)
		return err
	}

	sem := make(chan struct{}, workerPoolSize)
	var wg sync.WaitGroup

	for _, competition := range competitions {
		competitionsCached.Store(competition.ExternalId, competition)

		wg.Add(1)
		sem <- struct{}{}
		go func(cParam *entity.Competition) {
			defer wg.Done()
			defer func() { <-sem }()

			log.Printf("VOU BUSCAR OS TIMES DE: %s - ID: %s\n", cParam.Name, cParam.ExternalId)
			if err := s.fetchAndStoreTeams(ctx, cParam); err != nil {
				log.Printf("Erro ao buscar e armazenar times: %v", err)
			}
			teamsProcessedCount++
			log.Printf("(%v) TERMINEI DE BUSCAR OS TIMES DE: %s - ID: %s\n", teamsProcessedCount, cParam.Name, cParam.ExternalId)
		}(competition)
	}
	wg.Wait()

	var wgMatches sync.WaitGroup
	for _, competition := range competitions {
		wgMatches.Add(1)
		sem <- struct{}{}
		go func(cParam *entity.Competition) {
			defer wgMatches.Done()
			defer func() { <-sem }()

			log.Printf("VOU BUSCAR AS PARTIDAS DE: %s - ID: %s\n", cParam.Name, cParam.ExternalId)
			if err := s.fetchAndStoreMatches(ctx, cParam); err != nil {
				log.Printf("Erro ao buscar e armazenar partidas: %v", err)
			}
			matchesProcessedCount++
			log.Printf("(%v) TERMINEI DE BUSCAR AS PARTIDAS DE: %s - ID: %s\n", matchesProcessedCount, cParam.Name, cParam.ExternalId)
		}(competition)
	}
	wgMatches.Wait()

	s.printSummary()

	return nil
}

func (s *DataFetch) fetchAndStoreTeams(ctx context.Context, cEntity *entity.Competition) error {
	uriFragment := fmt.Sprintf("v4/competitions/%s/teams", cEntity.ExternalId)
	body, err := s.fetchFromExternalAPI(ctx, uriFragment)
	if err != nil {
		return err
	}

	var response dto.TeamResponseListExternalAPI
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	var teamEntities []*entity.Team
	for _, tResponse := range response.Teams {
		teamEntities = append(teamEntities, &entity.Team{
			UID:        uuid.New(),
			ExternalID: fmt.Sprintf("%d", tResponse.ID),
			Name:       tResponse.Name,
			FullName:   tResponse.FullName,
		})
	}

	teams, err := s.teamRepo.CreateOrUpdateInBatch(ctx, teamEntities)
	if err != nil {
		log.Printf("Erro ao salvar Times: %v", err)
		return err
	}

	for _, team := range teams {
		teamsCached.Store(team.ExternalID, team)
	}

	return nil
}

func (s *DataFetch) fetchAndStoreMatches(ctx context.Context, cEntity *entity.Competition) error {
	uriFragment := fmt.Sprintf("v4/competitions/%s/matches", cEntity.ExternalId)

	body, err := s.fetchFromExternalAPI(ctx, uriFragment)
	if err != nil {
		return err
	}

	var response dto.MatchResponseListExternalAPI
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	var cCachedEntity *entity.Competition
	var hTeamCachedEntity *entity.Team
	var aTeamCachedEntity *entity.Team

	var matchEntities []*entity.Match
	for _, mResponse := range response.Matches {

		if value, ok := competitionsCached.Load(fmt.Sprintf("%v", mResponse.Competition.ID)); ok {
			cCachedEntity = value.(*entity.Competition)
		}

		if value, ok := teamsCached.Load(fmt.Sprintf("%v", mResponse.HomeTeam.ID)); ok {
			hTeamCachedEntity = value.(*entity.Team)
		}

		if value, ok := teamsCached.Load(fmt.Sprintf("%v", mResponse.AwayTeam.ID)); ok {
			aTeamCachedEntity = value.(*entity.Team)
		}

		if cCachedEntity != nil && hTeamCachedEntity != nil && aTeamCachedEntity != nil {
			matchEntities = append(matchEntities, &entity.Match{
				UID:           uuid.New(),
				ExternalID:    fmt.Sprintf("%d", mResponse.ID),
				Round:         mResponse.Round,
				CompetitionID: cCachedEntity.ID,
				HomeTeam:      hTeamCachedEntity,
				AwayTeam:      aTeamCachedEntity,
				HomeTeamScore: nil,
				AwayTeamScore: nil,
			})
		}

	}
	_, err = s.matchRepo.CreateOrUpdateInBatch(ctx, matchEntities)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (s *DataFetch) fetchFromExternalAPI(ctx context.Context, uriFragment string) ([]byte, error) {
	apiURL := fmt.Sprintf("%s/%s", s.cfg.URL, uriFragment)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Auth-Token", s.cfg.Token)

	var body []byte
	fetchData := func() error {
		resp, err := s.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := resp.Header.Get("X-RequestCounter-Reset")
			if retryAfter != "" {
				waitTime, err := time.ParseDuration(retryAfter + "s")
				if err == nil {
					time.Sleep(waitTime)
				}
			}
			return fmt.Errorf("'429 TooManyRequests' retentar em : %ss", retryAfter)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("recebido status Não 200 : %d", resp.StatusCode)
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return nil
	}

	notify := func(err error, t time.Duration) {
		log.Printf("retentando `%s` após : %v\n", uriFragment, err)
	}

	retry := backoff.NewExponentialBackOff()
	retry.InitialInterval = retryInitialInterval * time.Second
	retry.MaxElapsedTime = retryMaxElapsedTime * retry.InitialInterval
	retry.RandomizationFactor = 0.75

	log.Printf("buscando dados de: %s \n", apiURL)
	err = backoff.RetryNotify(fetchData, retry, notify)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *DataFetch) printSummary() {
	competitionsLen := 0
	competitionsCached.Range(func(key, value interface{}) bool {
		competitionsLen++
		return true
	})

	teamsLen := 0
	teamsCached.Range(func(key, value interface{}) bool {
		teamsLen++
		return true
	})

	elapsedTime := time.Since(s.startTime)
	elapsedMinutes := int(elapsedTime.Minutes())
	elapsedSeconds := int(elapsedTime.Seconds()) % 60

	log.Println("\n\n------------------------------------------------------")
	log.Printf("COMPETIÇÕES INSERIDAS/ATUALIZADAS: %v", competitionsLen)
	log.Printf("TIMES INSERIDOS/ATUALIZADOS: %v", teamsLen)
	log.Printf("TEMPO DECORRIDO: %d:%02d minutos", elapsedMinutes, elapsedSeconds)
	log.Println("\n------------------------------------------------------")
}
