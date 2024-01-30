package main

import (
	"encoding/csv"
	"lolscout/api"
	"lolscout/data"
	"lolscout/tui"
	"os"
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

	matches, err := client.Get(dwx, queues).From(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		return
	}

	dwxFile, err := os.Create("dwx.csv")
	if err != nil {
		return
	}
	defer dwxFile.Close()

	marbeeFile, err := os.Create("marbee.csv")
	if err != nil {
		return
	}
	defer marbeeFile.Close()

	dwxWriter := csv.NewWriter(dwxFile)
	defer dwxWriter.Flush()

	marbeeWriter := csv.NewWriter(marbeeFile)
	defer marbeeWriter.Flush()

	header := []string{"champion", "level", "kills", "deaths", "assists", "kp", "cs", "cs/m", "win?", "matchType", "matchDurationMinutes"}

	dwxWriter.Write(header)
	marbeeWriter.Write(header)

	duos := data.FilterBySummoners(matches, []*lol.Summoner{dwx, marbee})

	for _, match := range duos {
		dwxStats := data.GetStats(match, dwx)

		dwxWriter.Write(dwxStats.Slice())

		marbeeStats := data.GetStats(match, marbee)

		marbeeWriter.Write(marbeeStats.Slice())

		println(tui.MatchParticipantModel{MatchParticipantStats: dwxStats}.View())
		println(tui.MatchParticipantModel{MatchParticipantStats: marbeeStats}.View())
		println()
	}
}
