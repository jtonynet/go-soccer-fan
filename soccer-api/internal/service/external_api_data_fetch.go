package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository"
)

type DataFetch struct {
	cfg             *config.ExternalApi
	areaRepo        repository.Area
	competitionRepo repository.Competition
	teamRepo        repository.Team
}

func NewDataFetchService(
	cfg *config.ExternalApi,
	areaRepo repository.Area,
	competitionRepo repository.Competition,
	teamRepo repository.Team) *DataFetch {
	return &DataFetch{
		cfg,
		areaRepo,
		competitionRepo,
		teamRepo,
	}
}

func (s *DataFetch) fetchFromExternalAPI(_ context.Context, uriFragment string) ([]byte, error) {
	apiURL := fmt.Sprintf("%s/%s", s.cfg.URL, uriFragment)
	fmt.Println(apiURL)

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
		fmt.Printf("retentando obter dados da api externa após erro: %v\n", err)
	}

	retry := backoff.NewExponentialBackOff()
	retry.InitialInterval = 5 * time.Second
	retry.MaxElapsedTime = 6 * retry.InitialInterval

	err = backoff.RetryNotify(fetchData, retry, notify)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *DataFetch) FetchAndStoreAreas(ctx context.Context) error {
	uriFragment := "/v4/areas"
	body, err := s.fetchFromExternalAPI(ctx, uriFragment)
	if err != nil {
		return err
	}

	var response dto.AreaResponseListExternalAPI
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	areaIDs := ""
	paramAreasCount := 0
	paramaAreasSize := 100
	var areasFilterParams []string
	var areaEntities []*entity.Area
	for idx, area := range response.Areas {
		externalID := fmt.Sprintf("%d", area.ID)
		areaIDs = fmt.Sprintf("%s,%s", areaIDs, externalID)

		areaEntities = append(areaEntities, &entity.Area{
			ExternalId:  externalID,
			Name:        area.Name,
			CountryCode: area.CountryCode,
		})

		if paramAreasCount == paramaAreasSize || idx+1 == len(response.Areas) {
			areasFilterParams = append(areasFilterParams, strings.TrimPrefix(areaIDs, ","))
			areaIDs = ""
			paramAreasCount = 0
			continue
		}

		paramAreasCount = paramAreasCount + 1

	}

	if err := s.areaRepo.CreateOrUpdateInBatch(ctx, areaEntities); err != nil {
		log.Printf("Erro ao salvar áreas: %v", err)
		return err
	}

	for _, areasParam := range areasFilterParams {
		go func(param string) {
			if err := s.FetchAndStoreCompetitions(context.Background(), param); err != nil {
				log.Printf("Erro ao buscar e armazenar competições: %v", err)
			}
		}(areasParam)
	}

	// fmt.Println(areasFilterParams)
	// s.FetchAndStoreCompetitions(context.Background(), "2032")

	return nil
}

func (s *DataFetch) FetchAndStoreCompetitions(ctx context.Context, areasFilterParam string) error {
	uriFragment := fmt.Sprintf("/v4/competitions/?areas=%s", areasFilterParam)
	body, err := s.fetchFromExternalAPI(ctx, uriFragment)
	if err != nil {
		return err
	}

	var response dto.CompetitionResponseListExternalAPI
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	var competitionEntities []*entity.Competition
	for _, competition := range response.Competitions {
		competitionEntities = append(competitionEntities, &entity.Competition{
			UID:        uuid.New(),
			ExternalId: fmt.Sprintf("%d", competition.ID),
			Name:       competition.Name,
			Season:     fmt.Sprintf("%v", time.Now().Year()),
		})
	}

	err = s.competitionRepo.CreateOrUpdateInBatch(ctx, competitionEntities)
	if err != nil {
		return err
	}

	for _, cEntity := range competitionEntities {
		go s.FetchAndStoreTeams(ctx, cEntity.ExternalId)
	}

	return nil
}

func (s *DataFetch) FetchAndStoreTeams(ctx context.Context, competitionExternalID string) error {
	uriFragment := fmt.Sprintf("v4/competitions/%s/teams", competitionExternalID)
	body, err := s.fetchFromExternalAPI(ctx, uriFragment)
	if err != nil {
		return err
	}

	var response dto.TeamResponseListExternalAPI
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	var entities []*entity.Team
	for _, team := range response.Teams {
		entities = append(entities, &entity.Team{
			UID:        uuid.New(),
			ExternalId: fmt.Sprintf("%d", team.ID),
			Name:       team.Name,
			FullName:   team.FullName,
		})
	}

	err = s.teamRepo.CreateOrUpdateInBatch(ctx, entities)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
