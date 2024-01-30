package data

import (
	"fmt"

	"github.com/KnutZuidema/golio/riot/lol"
)

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

func GetStats(match *lol.Match, summoner *lol.Summoner) MatchParticipantStats {
	teamKills := make(map[int]int)

	for _, participant := range match.Info.Participants {
		teamKills[participant.TeamID] += participant.Kills
	}

	durationMinutes := match.Info.GameDuration / 60

	for _, participant := range match.Info.Participants {
		if participant.PUUID == summoner.PUUID {
			var matchParticipant MatchParticipantStats

			matchParticipant.ChampionName = participant.ChampionName
			matchParticipant.Level = participant.ChampLevel
			matchParticipant.Kills = participant.Kills
			matchParticipant.Deaths = participant.Deaths
			matchParticipant.Assists = participant.Assists
			matchParticipant.KillParticipation = float64(participant.Kills+participant.Assists) / float64(teamKills[participant.TeamID])
			matchParticipant.CS = participant.TotalMinionsKilled
			matchParticipant.CSPerMinute = float64(participant.TotalMinionsKilled) / float64(durationMinutes)
			matchParticipant.Win = participant.Win
			matchParticipant.MatchType = lookupQueue(Queue(match.Info.QueueID))
			matchParticipant.DurationMinutes = durationMinutes

			return matchParticipant
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
