package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
}

func NewDataFetchService(
	cfg *config.ExternalApi,
	areaRepo repository.Area,
	competitionRepo repository.Competition) *DataFetch {
	return &DataFetch{
		cfg,
		areaRepo,
		competitionRepo,
	}
}

func (s *DataFetch) FetchAndStoreAreaData(ctx context.Context) error {
	apiURL := fmt.Sprintf("%s/v4/areas", s.cfg.URL)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Auth-Token", s.cfg.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response dto.AreaResponseListExternalAPI
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	areaIDs := ""
	paramAreasCount := 0
	paramaAreasSize := 20
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

	err = s.areaRepo.CreateOrUpdateInBatch(context.Background(), areaEntities)
	if err != nil {
		// TODO: tratamento de erro
		return err
	}

	for _, areasParam := range areasFilterParams {
		// println(fmt.Sprintf("%s/v4/competitions/?areas=%s", s.cfg.URL, areasParam))
		go s.FetchAndStoreCompetitionsData(context.Background(), areasParam)
	}

	return nil
}

func (s *DataFetch) FetchAndStoreCompetitionsData(ctx context.Context, areasFilterParam string) error {
	apiURL := fmt.Sprintf("%s/v4/competitions/?areas=%s", s.cfg.URL, areasFilterParam)
	println(apiURL)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Auth-Token", s.cfg.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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

	err = s.competitionRepo.CreateOrUpdateInBatch(context.Background(), competitionEntities)
	if err != nil {
		return err
	}

	return nil
}
