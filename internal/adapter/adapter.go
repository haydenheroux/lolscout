package adapter

import (
	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/haydenheroux/lolscout/internal/api/lol"
	"github.com/haydenheroux/lolscout/internal/metrics"
)

func GetMetrics(match *lol.Match, summoner *lol.Summoner) *metrics.MatchMetrics {
	teamDamage := make(map[int]int)
	teamKills := make(map[int]int)

	for _, participant := range match.Info.Participants {
		teamDamage[participant.TeamID] += participant.TotalDamageDealt
		teamKills[participant.TeamID] += participant.Kills
	}

	durationMinutes := float64(match.Info.GameDuration) / 60.0

	for _, participant := range match.Info.Participants {
		if participant.PUUID == summoner.PUUID {
			var metrics metrics.MatchMetrics

			metrics.Assists = participant.Assists
			metrics.CS = participant.TotalMinionsKilled + participant.NeutralMinionsKilled
			metrics.CSPerMinute = float64(metrics.CS) / float64(durationMinutes)
			metrics.ChampionName = participant.ChampionName
			metrics.ControlWardsPlaced = participant.DetectorWardsPlaced
			metrics.DamageDealt = participant.TotalDamageDealt
			metrics.DamageDealtPerMinute = float64(metrics.DamageDealt) / float64(durationMinutes)
			metrics.DamageDealtShare = float64(metrics.DamageDealt) / float64(teamDamage[participant.TeamID])
			metrics.Deaths = participant.Deaths
			metrics.DurationMinutes = durationMinutes
			metrics.KillParticipation = float64(participant.Kills+participant.Assists) / float64(teamKills[participant.TeamID])
			metrics.Kills = participant.Kills
			metrics.Level = participant.ChampLevel
			// TODO Refactor to include queue type?
			metrics.MatchType = api.QueueType(match.Info.QueueID).String()
			metrics.Position = participant.TeamPosition
			metrics.TurretsTaken = participant.TurretTakedowns
			metrics.WardsKilled = participant.WardsKilled
			metrics.WardsPlaced = participant.WardsPlaced
			metrics.Win = participant.Win

			return &metrics
		}
	}

	// TODO
	return &metrics.MatchMetrics{}
}
