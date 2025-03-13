package dto

type AreaFromExternalAPI struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
}

type AreaResponseExternalAPI struct {
	Count   int                    `json:"count"`
	Filters map[string]interface{} `json:"filters"`
	Areas   []AreaFromExternalAPI  `json:"areas"`
}
