package data

import (
	"fmt"
	"strconv"

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

func (stats MatchParticipantStats) Slice() []string {
	return []string{
		stats.ChampionName,
		strconv.Itoa(stats.Level),
		strconv.Itoa(stats.Kills),
		strconv.Itoa(stats.Deaths),
		strconv.Itoa(stats.Assists),
		strconv.FormatFloat(stats.KillParticipation, 'f', -1, 64),
		strconv.Itoa(stats.CS),
		strconv.FormatFloat(stats.CSPerMinute, 'f', -1, 64),
		strconv.FormatBool(stats.Win),
		stats.MatchType,
		strconv.Itoa(stats.DurationMinutes),
	}
}

func GetStats(match *lol.Match, summoner *lol.Summoner) MatchParticipantStats {
	teamKills := make(map[int]int)

	for _, participant := range match.Info.Participants {
		teamKills[participant.TeamID] += participant.Kills
	}

	durationMinutes := match.Info.GameDuration / 60

	for _, participant := range match.Info.Participants {
		if participant.PUUID == summoner.PUUID {
			var stats MatchParticipantStats

			stats.ChampionName = participant.ChampionName
			stats.Level = participant.ChampLevel
			stats.Kills = participant.Kills
			stats.Deaths = participant.Deaths
			stats.Assists = participant.Assists
			stats.KillParticipation = float64(participant.Kills+participant.Assists) / float64(teamKills[participant.TeamID])
			stats.CS = participant.TotalMinionsKilled + participant.NeutralMinionsKilled
			stats.CSPerMinute = float64(stats.CS) / float64(durationMinutes)
			stats.Win = participant.Win
			stats.MatchType = lookupQueue(Queue(match.Info.QueueID))
			stats.DurationMinutes = durationMinutes

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
