package main

import (
	"fmt"
	leagueApi "lolscout/api"
	"lolscout/tui"

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

	client := leagueApi.New(environment.RiotApiKey)

	summonerMatchParticipants, err := client.GetPlayer("dwx")

	if err != nil {
		log.Fatal(err)
	}

	for _, matchParticipant := range summonerMatchParticipants {
		fmt.Println(tui.RenderMatchParticipant(matchParticipant))
	}
}
