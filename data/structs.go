package data

import (
	"strconv"
)

type MatchParticipantMetrics struct {
	Assists              int
	CS                   int
	CSPerMinute          float64
	ChampionName         string
	ControlWardsPlaced   int
	DamageDealt          int
	DamageDealtPerMinute float64
	DamageDealtShare     float64
	Deaths               int
	DurationMinutes      float64
	KillParticipation    float64
	Kills                int
	Level                int
	MatchType            string
	Position             string
	TurretsTaken         int
	WardsKilled          int
	WardsPlaced          int
	Win                  bool
}

func formatBool(b bool) string {
	return strconv.FormatBool(b)
}

func formatInt(i int) string {
	return strconv.Itoa(i)
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func (stats MatchParticipantMetrics) Map() map[string]string {
	return map[string]string{
		"assists":         formatInt(stats.Assists),
		"cs":              formatInt(stats.CS),
		"cs/m":            formatFloat(stats.CSPerMinute),
		"champion":        stats.ChampionName,
		"controlWards":    formatInt(stats.ControlWardsPlaced),
		"dmg":             formatInt(stats.DamageDealt),
		"dmg/m":           formatFloat(stats.DamageDealtPerMinute),
		"dmg%":            formatFloat(stats.DamageDealtShare),
		"deaths":          formatInt(stats.Deaths),
		"durationMinutes": formatFloat(stats.DurationMinutes),
		"kp":              formatFloat(stats.KillParticipation),
		"kills":           formatInt(stats.Kills),
		"level":           formatInt(stats.Level),
		"matchType":       stats.MatchType,
		"posiiton":        stats.Position,
		"turrets":         formatInt(stats.TurretsTaken),
		"wardsKilled":     formatInt(stats.WardsKilled),
		"wardsPlaced":     formatInt(stats.WardsPlaced),
		"win":             formatBool(stats.Win),
	}
}

func (stats MatchParticipantMetrics) Header() []string {
	return []string{
		"assists",
		"cs",
		"cs/m",
		"champion",
		"controlWards",
		"dmg",
		"dmg/m",
		"dmg%",
		"deaths",
		"durationMinutes",
		"kp",
		"kills",
		"level",
		"matchType",
		"position",
		"turrets",
		"wardsKilled",
		"wardsPlaced",
		"win",
	}
}

func (stats MatchParticipantMetrics) Row() []string {
	return []string{
		formatInt(stats.Assists),
		formatInt(stats.CS),
		formatFloat(stats.CSPerMinute),
		stats.ChampionName,
		formatInt(stats.ControlWardsPlaced),
		formatInt(stats.DamageDealt),
		formatFloat(stats.DamageDealtPerMinute),
		formatFloat(stats.DamageDealtShare),
		formatInt(stats.Deaths),
		formatFloat(stats.DurationMinutes),
		formatFloat(stats.KillParticipation),
		formatInt(stats.Kills),
		formatInt(stats.Level),
		stats.MatchType,
		stats.Position,
		formatInt(stats.TurretsTaken),
		formatInt(stats.WardsKilled),
		formatInt(stats.WardsPlaced),
		formatBool(stats.Win),
	}
}
