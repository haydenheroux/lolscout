package main

import (
	"fmt"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
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

	client := golio.NewClient(environment.RiotApiKey,
		golio.WithRegion(api.RegionNorthAmerica),
		golio.WithLogger(log.New()))

	summoner, err := client.Riot.LoL.Summoner.GetByName("dwx")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s is a level %d summoner\n", summoner.Name, summoner.SummonerLevel)

	champion, _ := client.DataDragon.GetChampion("Irelia")
	mastery, err := client.Riot.LoL.ChampionMastery.Get(summoner.ID, champion.Key)
	if err != nil {
		fmt.Printf("%s has not played any games on %s\n", summoner.Name, champion.Name)
	} else {
		fmt.Printf("%s has mastery level %d with %d points on %s\n", summoner.Name, mastery.ChampionLevel,
			mastery.ChampionPoints, champion.Name)
	}
}
