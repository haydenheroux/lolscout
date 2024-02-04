package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"lolscout/api"
	"lolscout/data"
	"os"
	"strings"
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
				Usage: "Get data from recent matches",
				Subcommands: []*cli.Command{
					{
						Name:  "day",
						Usage: "Get data from the last day of matches",
						Action: func(c *cli.Context) error {
							return do(c, time.Now().AddDate(0, 0, -1))
						},
					},
					{
						Name:  "week",
						Usage: "Get data from the last week of matches",
						Action: func(c *cli.Context) error {
							return do(c, time.Now().AddDate(0, 0, -7))
						},
					},
					{
						Name:  "month",
						Usage: "Get data from last month of matches",
						Action: func(c *cli.Context) error {
							return do(c, time.Now().AddDate(0, -1, 0))
						},
					},
					{
						Name:  "year",
						Usage: "Get data from the last year of matches",
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

	fields := strings.Split(c.Args().First(), "#")

	if len(fields) != 2 {
		return errors.New("incorrect number of fields for Riot ID")
	}

	name := fields[0]
	tag := fields[1]

	summoner, err := client.TODO_SummonerByTag_TODO(name, tag)
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

	filename := fmt.Sprintf("%s#%s.csv", name, tag)

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
