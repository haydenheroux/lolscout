package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"lolscout/api"
	"lolscout/data"
	"os"
	"time"

	env "github.com/Netflix/go-env"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type Environment struct {
	RiotApiKey string `env:"RIOT_API_KEY,required=true"`
}

var environment Environment

func main() {
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "get recent matches",
				Subcommands: []*cli.Command{
					{
						Name:  "day",
						Usage: "gets the last day of matches",
						Action: func(c *cli.Context) error {
							return do(c, time.Now().AddDate(0, 0, -1))
						},
					},
					{
						Name:  "week",
						Usage: "gets the last week of matches",
						Action: func(c *cli.Context) error {
							return do(c, time.Now().AddDate(0, 0, -7))
						},
					},
					{
						Name:  "month",
						Usage: "gets the last month of matches",
						Action: func(c *cli.Context) error {
							return do(c, time.Now().AddDate(0, -1, 0))
						},
					},
					{
						Name:  "year",
						Usage: "gets the last year of matches",
						Action: func(c *cli.Context) error {
							return do(c, time.Now().AddDate(-1, 0, 0))
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func do(c *cli.Context, startTime time.Time) error {
	client := api.New(environment.RiotApiKey)

	puuid := c.Args().First()

	summoner, err := client.SummonerByPUUID(puuid)
	if err != nil {
		return err
	}

	queues := []data.Queue{data.Normal, data.Ranked, data.Clash}

	matches, err := client.Get(summoner, queues).From(startTime)
	if err != nil {
		return err
	}

	if len(matches) == 0 {
		return errors.New("summoner has no matches within the timeframe")
	}

	var stats []data.MatchParticipantStats

	for _, match := range matches {
		stats = append(stats, data.GetStats(match, summoner))
	}

	// TODO use name/tagline for filename
	filename := fmt.Sprintf("%s.csv", puuid)

	return writeCSV(filename, stats)
}

func writeCSV(name string, stats []data.MatchParticipantStats) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := false

	for _, s := range stats {
		if !header {
			err := writer.Write(s.Header())

			if err != nil {
				return err
			}

			header = true
		}

		err := writer.Write(s.Row())
		if err != nil {
			return err
		}
	}

	return nil
}
