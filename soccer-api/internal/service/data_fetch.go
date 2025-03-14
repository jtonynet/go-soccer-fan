package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository"
)

var retryInitialInterval = time.Duration(10)
var retryMaxElapsedTime = time.Duration(12)

var areasBatchSize = 100

var competitionsMap sync.Map
var teamsMap sync.Map

type DataFetch struct {
	cfg             *config.ExternalApi
	areaRepo        repository.Area
	competitionRepo repository.Competition
	teamRepo        repository.Team
	matchRepo       repository.Match
}

func NewDataFetchService(
	cfg *config.ExternalApi,
	areaRepo repository.Area,
	competitionRepo repository.Competition,
	teamRepo repository.Team,
	matchRepo repository.Match,
) *DataFetch {

	return &DataFetch{
		cfg,
		areaRepo,
		competitionRepo,
		teamRepo,
		matchRepo,
	}
}

func (s *DataFetch) fetchFromExternalAPI(_ context.Context, uriFragment string) ([]byte, error) {
	apiURL := fmt.Sprintf("%s/%s", s.cfg.URL, uriFragment)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Auth-Token", s.cfg.Token)

	var body []byte
	fetchData := func() error {
		client := &http.Client{}
		resp, err := client.Do(req)
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
			return fmt.Errorf("recebido status '429 TooManyRequests' retentar após : %ss", retryAfter)
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
		log.Printf("retentando obter dados da api externa após erro: %v\n", err)
	}

	retry := backoff.NewExponentialBackOff()
	retry.InitialInterval = retryInitialInterval * time.Second
	retry.MaxElapsedTime = retryMaxElapsedTime * retry.InitialInterval

	log.Printf("buscando dados de: %s \n", apiURL)
	err = backoff.RetryNotify(fetchData, retry, notify)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *DataFetch) FetchAndStore(ctx context.Context) error {
	uriFragment := "v4/areas"
	body, err := s.fetchFromExternalAPI(ctx, uriFragment)
	if err != nil {
		return err
	}

	var response dto.AreaResponseListExternalAPI
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	areaIDs := ""
	areasCount := 0
	var areasFilterParams []string
	var areaEntities []*entity.Area
	for idx, aResponse := range response.Areas {
		externalID := fmt.Sprintf("%d", aResponse.ID)
		areaIDs = fmt.Sprintf("%s,%s", areaIDs, externalID)

		areaEntities = append(areaEntities, &entity.Area{
			ExternalId:  externalID,
			Name:        aResponse.Name,
			CountryCode: aResponse.CountryCode,
		})

		if areasCount == areasBatchSize || idx+1 == len(response.Areas) {
			areasFilterParams = append(areasFilterParams, strings.TrimPrefix(areaIDs, ","))
			areaIDs = ""
			areasCount = 0
			continue
		}

		areasCount = areasCount + 1

	}

	if err := s.areaRepo.CreateOrUpdateInBatch(ctx, areaEntities); err != nil {
		log.Printf("erro ao salvar áreas: %v", err)
		return err
	}

	for _, areasParam := range areasFilterParams {
		go func(param string) {
			if err := s.fetchAndStoreCompetitions(context.Background(), param); err != nil {
				log.Printf("erro ao buscar e armazenar competições: %v", err)
			}
		}(areasParam)
	}

	return nil
}

func (s *DataFetch) fetchAndStoreCompetitions(ctx context.Context, areasFilterParam string) error {
	uriFragment := fmt.Sprintf("v4/competitions/?areas=%s", areasFilterParam)
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
		log.Printf("erro ao salvar competições: %v", err)
		return err
	}

	var wg sync.WaitGroup
	for _, competition := range competitions {
		competitionsMap.Store(competition.ExternalId, competition)

		wg.Add(1)
		go func(param *entity.Competition) {
			defer wg.Done()

			if err := s.fetchAndStoreTeams(ctx, competition); err != nil {
				log.Printf("erro ao buscar e armazenar times: %v", err)
				return
			}

			if err := s.fetchAndStoreMatches(ctx, competition); err != nil {
				log.Printf("erro ao buscar e armazenar partidas: %v", err)
				return
			}

		}(competition)
	}
	wg.Wait()

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
		log.Printf("erro ao salvar times: %v", err)
		return err
	}

	for _, team := range teams {
		teamsMap.Store(team.ExternalID, team)
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

	var cMappedEntity *entity.Competition
	var hTeamMappedEntity *entity.Team
	var aTeamMappedEntity *entity.Team

	var matchEntities []*entity.Match
	for _, mResponse := range response.Matches {

		if value, ok := competitionsMap.Load(fmt.Sprintf("%v", mResponse.Competition.ID)); ok {
			cMappedEntity = value.(*entity.Competition)
		}

		if value, ok := teamsMap.Load(fmt.Sprintf("%v", mResponse.HomeTeam.ID)); ok {
			hTeamMappedEntity = value.(*entity.Team)
		}

		if value, ok := teamsMap.Load(fmt.Sprintf("%v", mResponse.AwayTeam.ID)); ok {
			aTeamMappedEntity = value.(*entity.Team)
		}

		matchEntities = append(matchEntities, &entity.Match{
			UID:           uuid.New(),
			ExternalID:    fmt.Sprintf("%d", mResponse.ID),
			Round:         mResponse.Round,
			CompetitionID: cMappedEntity.ID,
			HomeTeam:      hTeamMappedEntity,
			AwayTeam:      aTeamMappedEntity,
			HomeTeamScore: nil,
			AwayTeamScore: nil,
		})

	}

	_, err = s.matchRepo.CreateOrUpdateInBatch(ctx, matchEntities)
	if err != nil {
		log.Printf("erro ao salvar times: %v", err)
		return err
	}

	return nil
}
