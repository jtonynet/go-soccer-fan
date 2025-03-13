package dto

type AreaResponseExternalAPI struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
}

type AreaResponseListExternalAPI struct {
	Count   int                       `json:"count"`
	Filters map[string]interface{}    `json:"filters"`
	Areas   []AreaResponseExternalAPI `json:"areas"`
}
