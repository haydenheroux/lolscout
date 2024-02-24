package main

import (
	"errors"
	"os"
	"time"

	env "github.com/Netflix/go-env"
	"github.com/haydenheroux/lolscout/internal/adapter"
	lol "github.com/haydenheroux/lolscout/internal/api/lol"
	playvs "github.com/haydenheroux/lolscout/internal/api/playvs"
	"github.com/haydenheroux/lolscout/internal/db"
	"github.com/haydenheroux/lolscout/internal/model"
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
				Name:  "lol",
				Usage: "Get League of Legends matches",
				Subcommands: []*cli.Command{
					{
						Name:  "get",
						Usage: "Get data from recent matches",
						Subcommands: []*cli.Command{
							{
								Name:  "day",
								Usage: "Get data from the last day of matches",
								Action: func(c *cli.Context) error {
									return scan(c.Args().First(), time.Now().AddDate(0, 0, -1))
								},
							},
							{
								Name:  "week",
								Usage: "Get data from the last week of matches",
								Action: func(c *cli.Context) error {
									return scan(c.Args().First(), time.Now().AddDate(0, 0, -7))
								},
							},
							{
								Name:  "month",
								Usage: "Get data from last month of matches",
								Action: func(c *cli.Context) error {
									return scan(c.Args().First(), time.Now().AddDate(0, -1, 0))
								},
							},
							{
								Name:  "year",
								Usage: "Get data from the last year of matches",
								Action: func(c *cli.Context) error {
									return scan(c.Args().First(), time.Now().AddDate(-1, 0, 0))
								},
							},
						},
					},
				},
			},
			{
				Name:  "playvs",
				Usage: "Get PlayVS data",
				Subcommands: []*cli.Command{
					{
						Name:  "teams",
						Usage: "Get PlayVS teams",
						Action: func(c *cli.Context) error {
							teamIds, err := playvs.CreateClient().Get().TeamIDs()

							if err != nil {
								return err
							}

							for _, teamId := range teamIds {
								println(teamId)
							}

							return nil
						},
					},
					{
						Name:  "players",
						Usage: "Get PlayVS players",
						Action: func(c *cli.Context) error {
							playvsClient := playvs.CreateClient()

							teamIds, err := playvsClient.Get().TeamIDs()

							if err != nil {
								return err
							}

							for _, teamId := range teamIds {

								playerNames, err := playvsClient.Get().Players(teamId)
								if err != nil {
									return err
								}

								for _, playerName := range playerNames {
									println(teamId, playerName)
								}
							}

							return nil
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

func scan(riotId string, startTime time.Time) error {
	client := lol.CreateClient(environment.RiotApiKey)

	var player model.Player

	summoner, err := client.TODO_SummonerByTag_TODO(riotId)
	if err != nil {
		return err
	}

	player.PUUID = summoner.PUUID

	queues := []lol.QueueType{lol.Queue.Normal, lol.Queue.Ranked, lol.Queue.Clash}

	matches, err := client.Get(summoner, queues).Since(startTime)
	if err != nil {
		return err
	}

	for _, match := range matches {
		metrics := adapter.GetMetrics(match, summoner)

		player.MatchMetrics = append(player.MatchMetrics, *metrics)
	}

	if len(matches) == 0 {
		return errors.New("summoner has no matches within the timeframe")
	}

	dbc, err := db.CreateClient("db.db")
	if err != nil {
		log.Fatal(err)
	}

	dbc.CreateOrUpdatePlayer(&player)

	return nil
}
