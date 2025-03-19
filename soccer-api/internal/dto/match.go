package dto

type MatchResponse struct {
	HomeTeamName string `json:"time_casa"`
	AwayTeamName string `json:"time_fora"`
	Score        string `json:"placar"`
}

type MatchResponseListPerRound struct {
	Round  int              `json:"rodada"`
	Matchs []*MatchResponse `json:"partidas,omitempty"`
}

type MatchResponseList struct {
	Rounds []*MatchResponseListPerRound `json:"rodadas"`
}

type TeamNested struct {
	ID int `json:"id"`
}

type CompetitionNested struct {
	ID int `json:"id"`
}

type Score struct {
	FullTime ScoreTime `json:"fullTime"`
	HalfTime ScoreTime `json:"halfTime"`
}

type ScoreTime struct {
	Home *int `json:"home"`
	Away *int `json:"away"`
}

type MatchResponseExternalAPI struct {
	ID          int               `json:"id"`
	Status      string            `json:"status"`
	Competition CompetitionNested `json:"competition"`
	HomeTeam    TeamNested        `json:"homeTeam"`
	AwayTeam    TeamNested        `json:"awayTeam"`
	Round       int               `json:"matchday"`
	Score       Score             `json:"score"`
}

type MatchResponseListExternalAPI struct {
	Filters map[string]interface{}     `json:"filters"`
	Matches []MatchResponseExternalAPI `json:"matches"`
}
