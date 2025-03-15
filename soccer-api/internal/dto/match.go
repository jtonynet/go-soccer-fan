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

type MatchResponseExternalAPI struct {
	ID            int               `json:"id"`
	Competition   CompetitionNested `json:"competition"`
	HomeTeam      TeamNested        `json:"homeTeam"`
	AwayTeam      TeamNested        `json:"awayTeam"`
	Round         int               `json:"matchday"`
	HomeTeamScore *int              `json:"home_team_score"`
	AwayTeamScore *int              `json:"away_team_score"`
}

type MatchResponseListExternalAPI struct {
	Filters map[string]interface{}     `json:"filters"`
	Matches []MatchResponseExternalAPI `json:"matches"`
}
