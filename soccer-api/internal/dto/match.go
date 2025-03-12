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
