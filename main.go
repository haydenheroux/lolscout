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

	if len(environment.RiotApiKey) == 0 {
		log.Fatal("RIOT_API_KEY missing from environment")
	}

	client := api.New(environment.RiotApiKey)

	dwx, err := client.Summoner("dwx")
	if err != nil {
		log.Fatal(err)
	}

	marbee, err := client.Summoner("marbee")
	if err != nil {
		log.Fatal(err)
	}

	queues := []data.Queue{data.Normal, data.Ranked, data.Clash}

	matches, err := client.Get(dwx, queues).From(time.Now().AddDate(0, -1, 0))
	if err != nil {
		log.Fatal(err)
	}

	var dwxStatsSlice, marbeeStatsSlice []data.MatchParticipantStats

	duos := data.FilterBySummoners(matches, []*lol.Summoner{dwx, marbee})

	for _, match := range duos {
		dwxStats := data.GetStats(match, dwx)
		dwxStatsSlice = append(dwxStatsSlice, dwxStats)

		marbeeStats := data.GetStats(match, marbee)
		marbeeStatsSlice = append(marbeeStatsSlice, marbeeStats)

		println(tui.MatchParticipantModel{MatchParticipantStats: dwxStats}.View())
		println(tui.MatchParticipantModel{MatchParticipantStats: marbeeStats}.View())
		println()
	}

	writeCSV("dwx.csv", dwxStatsSlice)
	writeCSV("marbee.csv", marbeeStatsSlice)
}

func writeCSV(name string, statsSlice []data.MatchParticipantStats) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := false

	for _, stats := range statsSlice {
		if !header {
			err := writer.Write(stats.Header())

			if err != nil {
				return err
			}

			header = true
		}

		err := writer.Write(stats.Row())
		if err != nil {
			return err
		}
	}

	return nil
}
