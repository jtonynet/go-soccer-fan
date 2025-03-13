package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jtonynet/go-soccer-fan/soccer-api/config"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/dto"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/repository"
)

type DataFetch struct {
	cfg      *config.ExternalApi
	areaRepo repository.Area
}

func NewDataFetchService(cfg *config.ExternalApi, areaRepo repository.Area) *DataFetch {
	return &DataFetch{
		cfg,
		areaRepo,
	}
}

func (s *DataFetch) FetchAndStoreData(ctx context.Context) error {
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

	var response dto.AreaResponseExternalAPI
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	var areas []*entity.Area
	for _, area := range response.Areas {
		areas = append(areas, &entity.Area{
			ExternalId:  fmt.Sprintf("%d", area.ID),
			Name:        area.Name,
			CountryCode: area.CountryCode,
		})
	}

	err = s.areaRepo.CreateOrUpdateInBatch(context.Background(), areas)
	if err != nil {
		return err
	}

	return nil
}
