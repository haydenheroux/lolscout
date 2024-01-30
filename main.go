package main

import (
	"fmt"
	"lolscout/api"
	"lolscout/data"
	"time"

	"github.com/KnutZuidema/golio/riot/lol"
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

	client := api.New(environment.RiotApiKey)

	dwx, err := client.Summoner("dwx")
	if err != nil {
		return
	}

	marbee, err := client.Summoner("marbee")
	if err != nil {
		return
	}

	queues := []int{400, 420, 700}

	matches, err := client.Get(dwx, queues).From(time.Now().AddDate(0, -1, 0))

	if err != nil {
		return
	}

	duos := data.FilterBySummoners(matches, []*lol.Summoner{dwx, marbee})

	dwxCount := 0
	marbeeCount := 0

	for _, match := range duos {
		dwxStats := data.GetStats(match, dwx)
		marbeeStats := data.GetStats(match, marbee)

		if dwxStats.CS > marbeeStats.CS {
			dwxCount += 1
		} else if marbeeStats.CS > dwxStats.CS {
			marbeeCount += 1
		}

		dwxString := "dwx " + summonerStatsToString(dwxStats)
		marbeeString := "marbee " + summonerStatsToString(marbeeStats)

		println(dwxString)
		println(marbeeString)
		println()
	}

	fmt.Printf("dwx total: %d marbee total: %d\n", dwxCount, marbeeCount)
}

func summonerStatsToString(data data.MatchParticipantStats) string {
	return fmt.Sprintf("(%s) cs: %d, cs/m: %.2f, kp: %.2f, won?: %v", data.ChampionName, data.CS, data.CSPerMinute, data.KillParticipation*100, data.Win)
}
