package stats

type MatchParticipantStats struct {
	ChampionName      string
	Level             int
	Kills             int
	Deaths            int
	Assists           int
	KillParticipation float64
	CS                int
	CSPerMinute       float64
	Win               bool
	MatchType         string
	DurationMinutes   int
}
