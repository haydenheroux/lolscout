package adapter

import (
	"time"

	"github.com/KnutZuidema/golio/riot/lol"
	riot "github.com/haydenheroux/lolscout/pkg/api/riot"
	"github.com/haydenheroux/lolscout/pkg/model"
)

func Team(id, name string, accounts []*riot.Account) *model.Team {
	var players []model.Player

	for _, account := range accounts {
		players = append(players, *Player(account))
	}

	return &model.Team{
		ID:      id,
		Name:    name,
		Players: players,
	}
}

func Player(account *riot.Account) *model.Player {
	return &model.Player{
		PUUID:    account.PUUID,
		GameName: account.GameName,
		TagLine:  account.TagLine,
	}
}

func MatchMetrics(match *lol.Match, summoner *lol.Summoner) *model.MatchMetrics {
	teamDamage := make(map[int]int)
	teamKills := make(map[int]int)

	for _, participant := range match.Info.Participants {
		teamDamage[participant.TeamID] += participant.TotalDamageDealt
		teamKills[participant.TeamID] += participant.Kills
	}

	durationMinutes := float64(match.Info.GameDuration) / 60.0

	for _, participant := range match.Info.Participants {
		if participant.PUUID == summoner.PUUID {
			var metrics model.MatchMetrics

			metrics.MatchID = match.Metadata.MatchID

			metrics.StartTime = time.UnixMilli(match.Info.GameStartTimestamp)

			metrics.Assists = participant.Assists
			metrics.CS = participant.TotalMinionsKilled + participant.NeutralMinionsKilled
			metrics.CSPerMinute = float64(metrics.CS) / float64(durationMinutes)

			metrics.Champion = model.Champion(participant.ChampionName)
			metrics.ControlWardsPlaced = participant.DetectorWardsPlaced
			metrics.DamageDealt = participant.TotalDamageDealt
			metrics.DamageDealtPerMinute = float64(metrics.DamageDealt) / float64(durationMinutes)
			metrics.DamageDealtShare = float64(metrics.DamageDealt) / float64(teamDamage[participant.TeamID])
			metrics.Deaths = participant.Deaths
			metrics.DurationMinutes = durationMinutes
			metrics.KillParticipation = float64(participant.Kills+participant.Assists) / float64(teamKills[participant.TeamID])
			metrics.Kills = participant.Kills
			metrics.Level = participant.ChampLevel
			metrics.MatchType = matchTypeOf(match)
			metrics.Position = positionOf(participant)
			metrics.TurretsTaken = participant.TurretTakedowns
			metrics.WardsKilled = participant.WardsKilled
			metrics.WardsPlaced = participant.WardsPlaced
			metrics.Win = participant.Win

			return &metrics
		}
	}

	return &model.MatchMetrics{}
}

// TODO
func matchTypeOf(match *lol.Match) model.MatchType {
	return model.MatchTypeSummonersRift
}

func positionOf(participant *lol.Participant) model.Position {
	switch participant.TeamPosition {
	case "TOP":
		return model.PositionTop
	case "JUNGLE":
		return model.PositionJungle
	case "MIDDLE":
		return model.PositionMiddle
	case "BOTTOM":
		return model.PositionBottom
	case "UTILITY":
		return model.PositionSupport
	}

	return model.Unknown
}
