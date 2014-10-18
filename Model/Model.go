package model

type UedBetData struct {
	TotalPages  int   `json:"tp"`
	LiveMatches Match `json:"i-ot"`
	PreMatches  Match `json:"n-ot"`
}

type Match struct {
	CategoryGroups []CategoryGroup `json:"egs"`
}

type CategoryGroup struct {
	Category   Category    `json:"c"`
	MatchInfos []MatchInfo `json:"es"`
}

type Category struct {
	Key  int64  `json:"k"`
	Name string `json:"n"`
}

type MatchInfo struct {
	Key         int      `json:"k"`
	Information []string `json:"i"`
	OddsInfo    OddsInfo `json:"o"`
}

type OddsInfo struct {
	HomeAway               []string `json:"1x2,omitempty"`
	HomeAwayFirstHalf      []string `json:"1x21st,omitempty"`
	AsianHandicap          []string `json:"ah,omitempty"`
	AsianHandicapFirstHalf []string `json:"ah1st,omitempty"`
	OverUnder              []string `json:"ou,omitempty"`
	OverUnderFirstHalf     []string `json:"ou1st,omitempty"`
	TotalGoal              []string `json:"tg,omitempty"`
	TotalGoalFirstHalf     []string `json:"tg1st,omitempty"`
}
