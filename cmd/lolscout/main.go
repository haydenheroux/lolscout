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
				Name:  "scan",
				Usage: "Scan recent matches",
				Subcommands: []*cli.Command{
					createLOLScanCommand("day", "Scan the last day of matches", 1),
					createLOLScanCommand("week", "Scan the last week of matches", 7),
					createLOLScanCommand("month", "Scan the last month of matches", 30),
					createLOLScanCommand("year", "Scan the last year of matches", 365),
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
				Name:  "teams",
				Usage: "PlayVS teams",
				Subcommands: []*cli.Command{
					{
						Name:  "init",
						Usage: "Initialize teams and players",
						Action: func(c *cli.Context) error {
							return initializePlayVSTeams()
						},
					},
					{
						Name:  "list",
						Usage: "List teams",
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
				},
			},
			{
				Name:  "team",
				Usage: "PlayVS team",
				Subcommands: []*cli.Command{
					{
						Name:  "info",
						Usage: "Show team information",
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

							for _, player := range team.Players {
								fmt.Printf("%s#%s\n", player.GameName, player.TagLine)
							}

							return nil
						},
					},
					{
						Name:  "scan",
						Usage: "Show team information",
						Action: func(c *cli.Context) error {
							dbc, err := db.CreateClient(environment.DatabaseName)
							if err != nil {
								return err
							}

							team, err := dbc.GetTeamByID(c.Args().First())

							if err != nil {
								return err
							}

							for _, player := range team.Players {
								// TODO Add variants for other time ranges
								monthAgo := time.Now().AddDate(0, 0, -30)

								err := scanLeagueOfLegendsMatches(player.GameName, player.TagLine, monthAgo)

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
}

func scanLeagueOfLegendsMatchesRiotId(riotId string, startTime time.Time) error {
	fields := strings.Split(riotId, "#")

	if len(fields) != 2 {
		return errors.New("bad riot id")
	}

	gameName := fields[0]
	tagLine := fields[1]

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

	player := adapter.Player(account)

	playerMatchIds, err := dbc.GetMatchIDsForPUUID(player.PUUID)

	summoner, err := lol.SummonerByPUUID(player.PUUID)

	queues := []lolApi.QueueType{lolApi.Queue.Normal, lolApi.Queue.Ranked, lolApi.Queue.Clash}

	matches, err := lol.Get(summoner, queues).Since(startTime)
	if err != nil {
		return err
	}

	log.Infof("got %d matches", len(matches))

	var matchMetrics []*model.MatchMetrics

	for _, match := range matches {
		if contains(playerMatchIds, match.Metadata.MatchID) {
			continue
		}

		metrics := adapter.MatchMetrics(match, summoner)

		matchMetrics = append(matchMetrics, metrics)

		player.PlayerMetrics = append(player.PlayerMetrics, *metrics)
	}

	log.Infof("saving %d matches (%d duplicates)", len(player.PlayerMetrics), len(matches)-len(player.PlayerMetrics))

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
			fields := strings.Split(displayName, "#")

			if len(fields) != 2 {
				log.Errorf("bad riot id format %s", displayName)
				continue
			}

			gameName := fields[0]
			tagLine := fields[1]

			account, err := riot.Get(gameName, tagLine).Account()
			if err != nil {
				log.Warnf("could not find riot id %s#%s", gameName, tagLine)
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
