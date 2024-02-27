package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	env "github.com/Netflix/go-env"
	"github.com/haydenheroux/lolscout/pkg/adapter"
	"github.com/haydenheroux/lolscout/pkg/analytics"
	lolApi "github.com/haydenheroux/lolscout/pkg/api/lol"
	playvsApi "github.com/haydenheroux/lolscout/pkg/api/playvs"
	riotApi "github.com/haydenheroux/lolscout/pkg/api/riot"
	"github.com/haydenheroux/lolscout/pkg/db"
	"github.com/haydenheroux/lolscout/pkg/model"
	"github.com/haydenheroux/lolscout/pkg/tui"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

type Environment struct {
	DatabaseName string `env:"DB_NAME,required=true"`
	RiotApiKey   string `env:"RIOT_API_KEY,required=true"`
}

var environment Environment

func main() {
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}

	app := createCLIApp()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func createCLIApp() *cli.App {
	app := &cli.App{
		Commands: []*cli.Command{
			createLOLCommand(),
			createPlayVSCommand(),
			{
				Name:  "thresholds",
				Usage: "prints per-role thresholds",
				Action: func(c *cli.Context) error {
					dbc, err := db.CreateClient(environment.DatabaseName)
					if err != nil {
						return err
					}

					thresholdsByPosition, err := dbc.GetPositionThresholds(0.6827)
					if err != nil {
						return err
					}

					first := true

					for position, thresholds := range thresholdsByPosition {
						if !first {
							fmt.Println()
						}

						first = false

						fmt.Println(position)
						fmt.Println(thresholds)
					}

					return nil
				},
			},
		},
	}
	return app
}

func createLOLCommand() *cli.Command {
	return &cli.Command{
		Name:  "lol",
		Usage: "League of Legends",
		Subcommands: []*cli.Command{
			{
				Name: "analyze",
				Action: func(c *cli.Context) error {
					riotId := c.Args().First()
					return analyzePlayer(riotId)
				},
			},
			{
				Name:  "scan",
				Usage: "scan recent matches",
				Subcommands: []*cli.Command{
					createLOLScanCommand("day", "scan the last day of matches", 1),
					createLOLScanCommand("week", "scan the last week of matches", 7),
					createLOLScanCommand("month", "scan the last month of matches", 30),
					createLOLScanCommand("year", "scan the last year of matches", 365),
				},
			},
		},
	}
}

func createLOLScanCommand(name, usage string, daysAgo int) *cli.Command {
	return &cli.Command{
		Name:  name,
		Usage: usage,
		Action: func(c *cli.Context) error {
			return scanLeagueOfLegendsMatchesRiotId(c.Args().First(), time.Now().AddDate(0, 0, -daysAgo))
		},
	}
}

func createPlayVSCommand() *cli.Command {
	return &cli.Command{
		Name:  "playvs",
		Usage: "PlayVS",
		Subcommands: []*cli.Command{
			{
				Name:  "analyze",
				Usage: "analyze a team's players",
				Action: func(c *cli.Context) error {
					dbc, err := db.CreateClient(environment.DatabaseName)
					if err != nil {
						return err
					}

					if c.Args().Len() != 1 {
						return errors.New("incorrect arguments")
					}

					team, err := dbc.GetTeamByID(c.Args().First())

					if err != nil {
						return err
					}

					for _, player := range team.Players {
						riotId := riotApi.Join(player.GameName, player.TagLine)

						fmt.Println(riotId)

						analyzePlayer(riotId)
					}

					return nil
				},
			},
			{
				Name:  "info",
				Usage: "display information for a team",
				Action: func(c *cli.Context) error {
					dbc, err := db.CreateClient(environment.DatabaseName)
					if err != nil {
						return err
					}

					team, err := dbc.GetTeamByID(c.Args().First())

					if err != nil {
						return err
					}

					fmt.Printf("%s: %s\n", team.Name, team.ID)
					fmt.Printf("has %d players\n", len(team.Players))

					for _, player := range team.Players {
						fmt.Printf("%s\n", riotApi.Join(player.GameName, player.TagLine))
					}

					return nil
				},
			},
			{
				Name:  "init",
				Usage: "initialize teams and players",
				Action: func(c *cli.Context) error {
					return initializePlayVSTeams()
				},
			},
			{
				Name:  "list",
				Usage: "list all teams",
				Action: func(c *cli.Context) error {
					dbc, err := db.CreateClient(environment.DatabaseName)
					if err != nil {
						return err
					}

					teams, err := dbc.GetAllTeams()

					for _, team := range teams {
						fmt.Printf("%s: %s\n", team.Name, team.ID)
					}

					return nil
				},
			},
			{
				Name:  "scan",
				Usage: "scan matches for a team",
				Subcommands: []*cli.Command{
					createPlayVSScanCommand("day", "scan the last day of matches", 1),
					createPlayVSScanCommand("week", "scan the last week of matches", 7),
					createPlayVSScanCommand("month", "scan the last month of matches", 30),
					createPlayVSScanCommand("year", "scan the last year of matches", 365),
				},
			},
		},
	}
}

func createPlayVSScanCommand(name, usage string, daysAgo int) *cli.Command {
	return &cli.Command{
		Name:  name,
		Usage: usage,
		Action: func(c *cli.Context) error {
			dbc, err := db.CreateClient(environment.DatabaseName)
			if err != nil {
				return err
			}

			teamId := c.Args().First()

			if len(teamId) == 0 {
				return errors.New("team id not specified")
			}

			team, err := dbc.GetTeamByID(teamId)

			if err != nil {
				return err
			}

			for _, player := range team.Players {
				err := scanLeagueOfLegendsMatches(player.GameName, player.TagLine, time.Now().AddDate(0, 0, -daysAgo))

				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}

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

	s14Start := time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)

	for _, metrics := range player.PlayerMetrics {
		if metrics.StartTime.After(s14Start) {
			s14PlayerMetrics = append(s14PlayerMetrics, metrics)
		}
	}

	thresholdsByPosition, err := dbc.GetPositionThresholds(0.6827)
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

	return nil
}
