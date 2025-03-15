package dto

type TeamResponseExternalAPI struct {
	ID       int    `json:"id"`
	Name     string `json:"shortName"`
	FullName string `json:"name"`
}

type TeamResponseListExternalAPI struct {
	Count   int                       `json:"count"`
	Filters map[string]interface{}    `json:"filters"`
	Teams   []TeamResponseExternalAPI `json:"teams"`
}
