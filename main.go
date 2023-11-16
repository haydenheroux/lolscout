package main

import (
	"fmt"
	"lolscout/tui"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
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

	client := golio.NewClient(environment.RiotApiKey,
		golio.WithRegion(api.RegionNorthAmerica),
		golio.WithLogger(log.New()))

	summoner, err := client.Riot.LoL.Summoner.GetByName("dwx")
	if err != nil {
		log.Fatal(err)
	}

	matchIds, err := client.Riot.LoL.Match.List(summoner.PUUID, 0, 20)
	if err != nil {
		log.Fatal(err)
	}

	var summonerMatchParticipants []tui.MatchParticipant

	for _, matchId := range matchIds {
		match, err := client.Riot.LoL.Match.Get(matchId)
		if err != nil {
			log.Fatal(err)
		}

		teamKills := make(map[int]int)

		for _, participant := range match.Info.Participants {
			teamKills[participant.TeamID] += participant.Kills
		}

		durationMinutes := match.Info.GameDuration / 60

		for _, participant := range match.Info.Participants {
			if participant.PUUID == summoner.PUUID {
				var matchParticipant tui.MatchParticipant

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

				summonerMatchParticipants = append(summonerMatchParticipants, matchParticipant)
			}
		}
	}

	for _, matchParticipant := range summonerMatchParticipants {
		fmt.Println(tui.RenderMatchParticipant(matchParticipant))
	}
}

func lookupQueue(queueId int) string {
	switch queueId {
	case 400:
		return "Normal"
	}
	return "TODO: Unimplemented"
}
