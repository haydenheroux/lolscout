package main

import (
	leagueApi "lolscout/api"
  
	env "github.com/Netflix/go-env"
	log "github.com/sirupsen/logrus"
)

type Environment struct {
	RiotApiKey string `env:"RIOT_API_KEY"`
}

func main() {
	var environment Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}

	client := leagueApi.New(environment.RiotApiKey)
  
	client.DoCS()
}

func matchToStats(match *lol.Match, summoner *lol.Summoner) stats.MatchParticipantStats {
	teamKills := make(map[int]int)

	for _, participant := range match.Info.Participants {
		teamKills[participant.TeamID] += participant.Kills
	}

	durationMinutes := match.Info.GameDuration / 60

	for _, participant := range match.Info.Participants {
		if participant.PUUID == summoner.PUUID {
			var matchParticipant stats.MatchParticipantStats

			matchParticipant.ChampionName = participant.ChampionName
			matchParticipant.Level = participant.ChampLevel
			matchParticipant.Kills = participant.Kills
			matchParticipant.Deaths = participant.Deaths
			matchParticipant.Assists = participant.Assists
			matchParticipant.KillParticipation = float64(participant.Kills+participant.Assists) / float64(teamKills[participant.TeamID])
			matchParticipant.CS = participant.TotalMinionsKilled
			matchParticipant.CSPerMinute = float64(participant.TotalMinionsKilled) / float64(durationMinutes)
			matchParticipant.Win = participant.Win
			matchParticipant.MatchType = lookupQueue(match.Info.QueueID)
			matchParticipant.DurationMinutes = durationMinutes

			return matchParticipant
		}
	}

	// TODO
	return stats.MatchParticipantStats{}
}

func lookupQueue(queueId int) string {
	switch queueId {
	case 400:
		return "Normal"
	}
	return fmt.Sprintf("TODO: %d", queueId)
}
