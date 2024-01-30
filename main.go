package main

import (
	"fmt"
	"lolscout/api"
	"lolscout/data"
	"lolscout/tui"
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

	queues := []data.Queue{data.Normal, data.Ranked, data.Clash}

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

		println(tui.MatchParticipantModel{MatchParticipantStats: dwxStats}.View())
		println(tui.MatchParticipantModel{MatchParticipantStats: marbeeStats}.View())
		println()
	}

	fmt.Printf("dwx total: %d marbee total: %d\n", dwxCount, marbeeCount)
}
