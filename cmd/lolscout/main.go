package main

import (
	"errors"
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
	"github.com/sirupsen/logrus"
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
							riot := riotApi.CreateClient(environment.RiotApiKey)
							playvs := playvsApi.CreateClient()

							dbc, err := db.CreateClient("db.db")
							if err != nil {
								return err
							}

							teams, err := playvs.GetRegion().Teams()

							if err != nil {
								return err
							}

							for _, team := range teams {
								roster, err := playvs.GetRegion().GetTeam(team).Roster()
								if err != nil {
									return err
								}

								var accounts []*riotApi.Account

								for _, displayName := range roster.DisplayNames {
									fields := strings.Split(displayName, "#")

									if len(fields) != 2 {
										logrus.Errorf("bad riot id format %s", displayName)
										continue
									}

									gameName := fields[0]
									tagLine := fields[1]

									account, err := riot.Get(gameName, tagLine).Account()
									if err != nil {
										logrus.Warnf("could not find riot id %s#%s", gameName, tagLine)
										logrus.Infof("reason: %v", err)
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

	logrus.Infof("got %d matches", len(matches))

	var matchMetrics []*model.MatchMetrics

	for _, match := range matches {
		// Skip don't append if already stored
		if contains(playerMatchIds, match.Metadata.MatchID) {
			continue
		}

		metrics := adapter.MatchMetrics(match, summoner)

		matchMetrics = append(matchMetrics, metrics)

		player.PlayerMetrics = append(player.PlayerMetrics, *metrics)
	}

	logrus.Infof("saving %d matches (%d duplicates)", len(player.PlayerMetrics), len(matches)-len(player.PlayerMetrics))

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
