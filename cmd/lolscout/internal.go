package main

import (
	"errors"
	"fmt"

	"github.com/haydenheroux/lolscout/pkg/adapter"
	"github.com/haydenheroux/lolscout/pkg/analytics"
	"github.com/haydenheroux/lolscout/pkg/db"
	"github.com/haydenheroux/lolscout/pkg/model"
	"github.com/haydenheroux/lolscout/pkg/tui"
	"gorm.io/gorm"
)

func scanLeagueOfLegendsMatchesRiotId(riotId string, startTime time.Time) error {
	gameName, tagLine, err := riotApi.Split(riotId)

	if err != nil {
		return err
	}

	return scanLeagueOfLegendsMatches(gameName, tagLine, startTime)
}

func scanLeagueOfLegendsMatches(gameName, tagLine string, startTime time.Time) error {
	riot := riotApi.CreateClient(environment.RiotApiKey)
	lol := lolApi.CreateClient(environment.RiotApiKey)

	account, err := riot.Get(gameName, tagLine).Account()
	if err != nil {
		return err
	}

	dbc, err := db.CreateClient(environment.DatabaseName)
	if err != nil {
		return err
	}

	puuid := account.PUUID

	var player *model.Player

	p, err := dbc.GetPlayerByPUUID(puuid)
	if err == nil {
		player = p
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		player = &model.Player{
			PUUID:         puuid,
			GameName:      gameName,
			TagLine:       tagLine,
			TeamID:        nil,
			PlayerMetrics: make([]model.MatchMetrics, 0),
		}
	} else {
		return err
	}

	summoner, err := lol.SummonerByPUUID(puuid)
	if err != nil {
		return err
	}

	queues := []lolApi.QueueType{lolApi.Queue.Normal, lolApi.Queue.Ranked, lolApi.Queue.Clash}

	matches, err := lol.Get(summoner, queues).Since(startTime)
	if err != nil {
		return err
	}

	log.Infof("got %d matches", len(matches))

	var matchMetrics []*model.MatchMetrics

	for _, match := range matches {
		metrics := adapter.MatchMetrics(match, summoner)

		matchMetrics = append(matchMetrics, metrics)
	}

	scanned := len(matchMetrics)
	appended := player.AppendMatchMetrics(matchMetrics)

	log.Infof("saving %d matches (%d duplicates)", appended, scanned-appended)

	err = dbc.CreateOrUpdatePlayer(player)
	if err != nil {
		return err
	}

	return nil
}

func initializePlayVSTeams() error {
	riot := riotApi.CreateClient(environment.RiotApiKey)
	playvs := playvsApi.CreateClient()

	region := playvs.GetRegion(playvsApi.EasternRegion)

	dbc, err := db.CreateClient(environment.DatabaseName)
	if err != nil {
		return err
	}

	teams, err := region.Teams()

	if err != nil {
		return err
	}

	for _, team := range teams {
		displayNames, err := region.Get(team).DisplayNames()
		if err != nil {
			return err
		}

		var accounts []*riotApi.Account

		for _, displayName := range displayNames {
			gameName, tagLine, err := riotApi.Split(displayName)

			if err != nil {
				log.Warnf("%v %s", err, displayName)
				continue
			}

			account, err := riot.Get(gameName, tagLine).Account()
			if err != nil {
				log.Warnf("could not find riot id %s", displayName)
				log.Infof("reason: %v", err)
				continue
			}

			accounts = append(accounts, account)
		}

		err = dbc.CreateOrUpdateTeam(adapter.Team(team.ID, team.Name, accounts))

		if err != nil {
			return err
		}
	}

	return nil
}

func analyzePlayer(riotId string) error {
	riot := riotApi.CreateClient(environment.RiotApiKey)

	name, tag, err := riotApi.Split(riotId)
	if err != nil {
		return err
	}

	account, err := riot.Get(name, tag).Account()
	if err != nil {
		return err
	}

	dbc, err := db.CreateClient(environment.DatabaseName)
	if err != nil {
		return err
	}

	player, err := dbc.GetPlayerByPUUID(account.PUUID)
	if err != nil {
		return err
	}

	var s14PlayerMetrics []model.MatchMetrics

	for _, metrics := range player.PlayerMetrics {
		if metrics.StartTime.After(time.Season14()) {
			s14PlayerMetrics = append(s14PlayerMetrics, metrics)
		}
	}

	thresholdsByPosition, err := dbc.GetPositionThresholds(0.6827)
	if err != nil {
		return err
	}

	thresholdsByChampion, err := dbc.GetChampionThresholds(0.6827)
	if err != nil {
		return err
	}

	analyticsByPosition := analytics.AnalyzeByPosition(s14PlayerMetrics)

	first := true

	for position, analytics := range analyticsByPosition {
		if analytics.Size > 2 {
			if !first {
				fmt.Println()
			}

			first = false

			fmt.Println(position)

			a := tui.Analytics{Analytics: analytics, Thresholds: thresholdsByPosition[position]}

			fmt.Println(a.View())
		}
	}

	analyticsByChampion := analytics.AnalyzeByChampion(s14PlayerMetrics)

	first = true

	for champion, analytics := range analyticsByChampion {
		if analytics.Size > 2 {
			if !first {
				fmt.Println()
			}

			first = false

			fmt.Println(champion)

			a := tui.Analytics{Analytics: analytics, Thresholds: thresholdsByChampion[champion]}

			fmt.Println(a.View())
		}
	}

	return nil
}
