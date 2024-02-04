package data

import (
	"fmt"
	"strconv"

	"github.com/KnutZuidema/golio/riot/lol"
)

type MatchParticipantStats struct {
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

func (stats MatchParticipantStats) Map() map[string]string {
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
		"turrets":         formatInt(stats.TurretsTaken),
		"wardsKilled":     formatInt(stats.WardsKilled),
		"wardsPlaced":     formatInt(stats.WardsPlaced),
		"win":             formatBool(stats.Win),
	}
}

func (stats MatchParticipantStats) Header() []string {
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
		"turrets",
		"wardsKilled",
		"wardsPlaced",
		"win",
	}
}

func (stats MatchParticipantStats) Row() []string {
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
		formatInt(stats.TurretsTaken),
		formatInt(stats.WardsKilled),
		formatInt(stats.WardsPlaced),
		formatBool(stats.Win),
	}
}

func GetStats(match *lol.Match, summoner *lol.Summoner) MatchParticipantStats {
	teamDamage := make(map[int]int)
	teamKills := make(map[int]int)

	for _, participant := range match.Info.Participants {
		teamDamage[participant.TeamID] += participant.TotalDamageDealt
		teamKills[participant.TeamID] += participant.Kills
	}

	durationMinutes := float64(match.Info.GameDuration) / 60.0

	for _, participant := range match.Info.Participants {
		if participant.PUUID == summoner.PUUID {
			var stats MatchParticipantStats

			stats.Assists = participant.Assists
			stats.CS = participant.TotalMinionsKilled + participant.NeutralMinionsKilled
			stats.CSPerMinute = float64(stats.CS) / float64(durationMinutes)
			stats.ChampionName = participant.ChampionName
			stats.ControlWardsPlaced = participant.DetectorWardsPlaced
			stats.DamageDealt = participant.TotalDamageDealt
			stats.DamageDealtPerMinute = float64(stats.DamageDealt) / float64(durationMinutes)
			stats.DamageDealtShare = float64(stats.DamageDealt) / float64(teamDamage[participant.TeamID])
			stats.Deaths = participant.Deaths
			stats.DurationMinutes = durationMinutes
			stats.KillParticipation = float64(participant.Kills+participant.Assists) / float64(teamKills[participant.TeamID])
			stats.Kills = participant.Kills
			stats.Level = participant.ChampLevel
			stats.MatchType = lookupQueue(Queue(match.Info.QueueID))
			stats.TurretsTaken = participant.TurretTakedowns
			stats.WardsKilled = participant.WardsKilled
			stats.WardsPlaced = participant.WardsPlaced
			stats.Win = participant.Win

			return stats
		}
	}

	// TODO
	return MatchParticipantStats{}
}

type Queue int

const (
	Normal Queue = 400
	Ranked Queue = 420
	Clash  Queue = 700
)

func lookupQueue(queue Queue) string {
	switch queue {
	case Normal:
		return "Normal"
	case Ranked:
		return "Ranked"
	case Clash:
		return "Clash"
	}
	return fmt.Sprintf("TODO: %d", queue)
}
