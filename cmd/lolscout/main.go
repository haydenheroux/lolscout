package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	env "github.com/Netflix/go-env"
	"github.com/haydenheroux/lolscout/internal/adapter"
	lolApi "github.com/haydenheroux/lolscout/internal/api/lol"
	playvsApi "github.com/haydenheroux/lolscout/internal/api/playvs"
	riotApi "github.com/haydenheroux/lolscout/internal/api/riot"
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
							teamIds, err := playvsApi.CreateClient().Get().TeamIDs()

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
							playvsClient := playvsApi.CreateClient()

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
	riot := riotApi.CreateClient(environment.RiotApiKey)
	lol := lolApi.CreateClient(environment.RiotApiKey)

	fields := strings.Split(riotId, "#")

	if len(fields) != 2 {
		return errors.New("bad riot id")
	}

	gameName := fields[0]
	tagLine := fields[1]

	account, err := riot.Get(gameName, tagLine).Account()
	if err != nil {
		return err
	}

	dbc, err := db.CreateClient("db.db")
	if err != nil {
		return err
	}

	player := adapter.Player(account)

	playerMatchIds, err := dbc.GetMatchIDsForPUUID(player.PUUID)

	summoner, err := lol.SummonerByPUUID(player.PUUID)

	queues := []lolApi.QueueType{lolApi.Queue.Normal, lolApi.Queue.Ranked, lolApi.Queue.Clash}

	matches, err := lol.Get(summoner, queues).Since(startTime)
	if err != nil {
		return err
	}

	fmt.Printf("got %d matches\n", len(matches))

	var matchMetrics []*model.MatchMetrics

	for _, match := range matches {
		// Skip don't append if already stored
		if contains(playerMatchIds, match.Metadata.MatchID) {
			continue
		}

		metrics := adapter.GetMetrics(match, summoner)

		matchMetrics = append(matchMetrics, metrics)

		player.PlayerMetrics = append(player.PlayerMetrics, *metrics)
	}

	fmt.Printf("saving %d matches (%d duplicates)\n", len(player.PlayerMetrics), len(matches)-len(player.PlayerMetrics))

	err = dbc.CreateOrUpdatePlayer(player)
	if err != nil {
		return err
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
