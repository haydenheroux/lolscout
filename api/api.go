package api

import (
	"fmt"
	"time"

	"lolscout/data"

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

func (api LeagueAPI) GetPlayer(name string) ([]data.MatchParticipantStats, error) {
	summoner, err := api.client.Riot.LoL.Summoner.GetByName(name)
	if err != nil {
		return []data.MatchParticipantStats{}, err
	}

	matchIds, err := api.client.Riot.LoL.Match.List(summoner.PUUID, 0, 20)
	if err != nil {
		return []data.MatchParticipantStats{}, err
	}

	var summonerMatchParticipants []data.MatchParticipantStats

	for _, matchId := range matchIds {
		match, err := api.client.Riot.LoL.Match.Get(matchId)
		if err != nil {
			return []data.MatchParticipantStats{}, err
		}

		summonerMatchParticipants = append(summonerMatchParticipants, transformMatch(match, summoner))
	}

	return summonerMatchParticipants, nil
}

func (api LeagueAPI) DoCS() {
	dwx, err := api.client.Riot.LoL.Summoner.GetByName("dwx")
	if err != nil {
		return
	}

	marbee, err := api.client.Riot.LoL.Summoner.GetByName("marbee")
	if err != nil {
		return
	}

	matches, err := api.getMonthMatches(dwx)

	if err != nil {
		return
	}

	duos, err := api.filterMatchesByParticipants(matches, []*lol.Summoner{dwx, marbee})

	if err != nil {
		return
	}

	var duoSR []lol.Match

	for _, match := range duos {
		queue := match.Info.QueueID

		ranked := 420
		normal := 400
		clash := 700

		if queue == ranked || queue == normal || queue == clash {
			duoSR = append(duoSR, match)
		}
	}

	dwxCount := 0
	marbeeCount := 0

	for _, match := range duoSR {
		dwxData := transformMatch(&match, dwx)
		marbeeData := transformMatch(&match, marbee)

		dwxCS := dwxData.CS
		marbeeCS := marbeeData.CS

		dwxString := "dwx " + summonerStatsToString(dwxData)
		marbeeString := "marbee " + summonerStatsToString(marbeeData)

		println(dwxString)
		println(marbeeString)
		println()

		if dwxCS > marbeeCS {
			dwxCount += 1
		} else if marbeeCS > dwxCS {
			marbeeCount += 1
		}
	}

	fmt.Printf("dwx total: %d marbee total: %d\n", dwxCount, marbeeCount)
}

func summonerStatsToString(data data.MatchParticipantStats) string {
	return fmt.Sprintf("(%s) cs: %d, cs/m: %.2f, kp: %.2f, won?: %v", data.ChampionName, data.CS, data.CSPerMinute, data.KillParticipation*100, data.Win)
}

func (api LeagueAPI) getMonthMatches(summoner *lol.Summoner) ([]lol.Match, error) {
	// monthAgo := time.Now().AddDate().AddDate(0, -1, 0)
	monthAgo := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)

	predicate := func(match lol.Match) bool {
		sec := match.Info.GameStartTimestamp / 1000

		matchTime := time.Unix(sec, 0)

		println("got match")

		return matchTime.After(monthAgo)
	}

	return api.getMatchesUntil(summoner, predicate)
}

func (api LeagueAPI) getMatchesUntil(summoner *lol.Summoner, predicate func(lol.Match) bool) ([]lol.Match, error) {
	var matches []lol.Match

	for result := range api.client.Riot.LoL.Match.ListStream(summoner.PUUID) {
		if result.Error != nil {
			break
		}

		match, err := api.client.Riot.LoL.Match.Get(result.MatchID)
		if err != nil {
			break
		}

		if !predicate(*match) {
			break
		}

		matches = append(matches, *match)
	}

	return matches, nil
}

func (api LeagueAPI) filterMatchesByParticipants(matches []lol.Match, summoners []*lol.Summoner) ([]lol.Match, error) {
	var result []lol.Match

	for _, match := range matches {
		if hasParticipants(match, summoners) {
			result = append(result, match)
		}
	}

	return result, nil
}

func hasParticipants(match lol.Match, summoners []*lol.Summoner) bool {
	for _, summoner := range summoners {
		participated := false

		for _, participantPUUID := range match.Metadata.Participants {
			if summoner.PUUID == participantPUUID {
				participated = true
			}
		}

		if !participated {
			return false
		}
	}

	return true
}

func (api LeagueAPI) getPUUIDS(names []string) ([]string, error) {
	var puuids []string

	for _, name := range names {
		summoner, err := api.client.Riot.LoL.Summoner.GetByName(name)

		if err != nil {
			return []string{}, err
		}

		puuids = append(puuids, summoner.PUUID)
	}

	return puuids, nil
}

func transformMatch(match *lol.Match, summoner *lol.Summoner) data.MatchParticipantStats {
	teamKills := make(map[int]int)

	for _, participant := range match.Info.Participants {
		teamKills[participant.TeamID] += participant.Kills
	}

	durationMinutes := match.Info.GameDuration / 60

	for _, participant := range match.Info.Participants {
		if participant.PUUID == summoner.PUUID {
			var matchParticipant data.MatchParticipantStats

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
	return data.MatchParticipantStats{}
}

func lookupQueue(queueId int) string {
	switch queueId {
	case 400:
		return "Normal"
	}
	return fmt.Sprintf("TODO: %d", queueId)
}
