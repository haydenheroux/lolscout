package api

import (
	"lolscout/structs"

	"github.com/KnutZuidema/golio"

	"github.com/KnutZuidema/golio/api"
	"github.com/KnutZuidema/golio/riot/lol"
	log "github.com/sirupsen/logrus"
)

type LeagueAPI struct {
	client *golio.Client
}

func New(key string) LeagueAPI {
	return LeagueAPI{
		client: golio.NewClient(key,
			golio.WithRegion(api.RegionNorthAmerica),
			golio.WithLogger(log.New())),
	}
}

func (api LeagueAPI) GetPlayer(name string) ([]structs.MatchParticipant, error) {
	summoner, err := api.client.Riot.LoL.Summoner.GetByName(name)
	if err != nil {
		return []structs.MatchParticipant{}, err
	}

	matchIds, err := api.client.Riot.LoL.Match.List(summoner.PUUID, 0, 20)
	if err != nil {
		return []structs.MatchParticipant{}, err
	}

	var summonerMatchParticipants []structs.MatchParticipant

	for _, matchId := range matchIds {
		match, err := api.client.Riot.LoL.Match.Get(matchId)
		if err != nil {
			return []structs.MatchParticipant{}, err
		}

		summonerMatchParticipants = append(summonerMatchParticipants, transformMatch(match, summoner))
	}

	return summonerMatchParticipants, nil
}

func transformMatch(match *lol.Match, summoner *lol.Summoner) structs.MatchParticipant {
	var matchParticipant structs.MatchParticipant

	teamKills := make(map[int]int)

	for _, participant := range match.Info.Participants {
		teamKills[participant.TeamID] += participant.Kills
	}

	durationMinutes := match.Info.GameDuration / 60

	for _, participant := range match.Info.Participants {
		if participant.PUUID == summoner.PUUID {
			matchParticipant.ChampionName = participant.ChampionName
			matchParticipant.Level = participant.ChampLevel
			matchParticipant.Kills = participant.Kills
			matchParticipant.Deaths = participant.Deaths
			matchParticipant.Assists = participant.Assists
			matchParticipant.KillParticipation = float64(participant.Kills+participant.Assists) / float64(teamKills[participant.TeamID])
			matchParticipant.CS = participant.TotalMinionsKilled
			matchParticipant.CSPerMinute = float64(participant.TotalMinionsKilled) / float64(durationMinutes)
			matchParticipant.Win = participant.Win
			break
		}
	}

	return matchParticipant
}
