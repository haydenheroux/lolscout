package main

import (
	"fmt"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	env "github.com/Netflix/go-env"
	"github.com/charmbracelet/lipgloss"
	log "github.com/sirupsen/logrus"
)

type Environment struct {
	RiotApiKey string `env:"RIOT_API_KEY"`
}

type MatchParticipant struct {
	ChampionName      string
	Level             int
	Kills             int
	Deaths            int
	Assists           int
	KillParticipation float64
	CS                int
	CSPerMinute       float64
	Win               bool
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

	matchIds, err := client.Riot.LoL.Match.List(summoner.PUUID, 0, 20)
	if err != nil {
		log.Fatal(err)
	}

	var summonerMatchParticipants []MatchParticipant

	for _, matchId := range matchIds {
		match, err := client.Riot.LoL.Match.Get(matchId)
		if err != nil {
			log.Fatal(err)
		}

		teamKills := make(map[int]int)

		for _, participant := range match.Info.Participants {
			teamKills[participant.TeamID] += participant.Kills
		}

		durationMinutes := match.Info.GameDuration / 60

		for _, participant := range match.Info.Participants {
			if participant.PUUID == summoner.PUUID {
				var matchParticipant MatchParticipant

				matchParticipant.ChampionName = participant.ChampionName
				matchParticipant.Level = participant.ChampLevel
				matchParticipant.Kills = participant.Kills
				matchParticipant.Deaths = participant.Deaths
				matchParticipant.Assists = participant.Assists
				matchParticipant.KillParticipation = float64(participant.Kills+participant.Assists) / float64(teamKills[participant.TeamID])
				matchParticipant.CS = participant.TotalMinionsKilled
				matchParticipant.CSPerMinute = float64(participant.TotalMinionsKilled) / float64(durationMinutes)
				matchParticipant.Win = participant.Win

				summonerMatchParticipants = append(summonerMatchParticipants, matchParticipant)
			}
		}
	}

	for _, matchParticipant := range summonerMatchParticipants {
		var backgroundColor, borderColor lipgloss.Color

		if matchParticipant.Win {
			backgroundColor = lipgloss.Color("#28344E")
			borderColor = lipgloss.Color("#5383E8")
		} else {
			backgroundColor = lipgloss.Color("#59343B")
			borderColor = lipgloss.Color("#E84057")
		}

		// TODO
		championSectionWidth := 14

		championName := lipgloss.NewStyle().Background(backgroundColor).Bold(true).Width(championSectionWidth)

		renderedChampionName := championName.Render(matchParticipant.ChampionName)

		levelString := fmt.Sprintf("Lvl. %d", matchParticipant.Level)

		level := lipgloss.NewStyle().Background(backgroundColor).PaddingRight(1).Width(championSectionWidth)

		renderedLevel := level.Render(levelString)

		renderedChampionSection := lipgloss.JoinVertical(lipgloss.Left, renderedChampionName, renderedLevel)

		kdaValue := float64(matchParticipant.Kills+matchParticipant.Assists) / float64(matchParticipant.Deaths)

		whiteKda := lipgloss.NewStyle().Background(backgroundColor)

		goldKda := whiteKda.Copy().Foreground(lipgloss.Color("#FF8200"))
		blueKda := goldKda.Copy().Foreground(lipgloss.Color("#0093FF"))
		greenKda := goldKda.Copy().Foreground(lipgloss.Color("#00BBA3"))

		kdaTextStyle := whiteKda

		if kdaValue >= 5 {
			kdaTextStyle = goldKda
		} else if kdaValue >= 4 {
			kdaTextStyle = blueKda
		} else if kdaValue >= 3 {
			kdaTextStyle = greenKda
		}

		kdaTextStyle = kdaTextStyle.PaddingRight(1)

		kdaText := fmt.Sprintf("%.2f:1 KDA", kdaValue)

		renderedKdaText := kdaTextStyle.Render(kdaText)

		killParticipationString := fmt.Sprintf("(%.0f%% KP)", matchParticipant.KillParticipation*100)

		killParticipation := lipgloss.NewStyle().Background(backgroundColor).PaddingRight(1)

		if matchParticipant.KillParticipation >= 0.5 {
			killParticipation = killParticipation.Foreground(lipgloss.Color("#E84057"))
		}

		renderedKillParticipation := killParticipation.Render(killParticipationString)

		renderedKdaBottomSection := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedKdaText, renderedKillParticipation)

		kdaString := fmt.Sprintf("%d/%d/%d", matchParticipant.Kills, matchParticipant.Deaths, matchParticipant.Assists)

		kda := lipgloss.NewStyle().Background(backgroundColor).Width(lipgloss.Width(renderedKdaBottomSection))

		renderedKda := kda.Render(kdaString)

		// TODO Center renderedKda above bottom section
		renderedKdaSection := lipgloss.JoinVertical(lipgloss.Left, renderedKda, renderedKdaBottomSection)

		csString := fmt.Sprintf("%d CS", matchParticipant.CS)

		csPerMinuteString := fmt.Sprintf("%.1f CS/M", matchParticipant.CSPerMinute)

		csPerMinute := lipgloss.NewStyle().Background(backgroundColor)

		if matchParticipant.CSPerMinute >= 8.0 {
			csPerMinute = csPerMinute.Foreground(lipgloss.Color("#E84057"))
		}

		renderedCsPerMinute := csPerMinute.Render(csPerMinuteString)

		cs := lipgloss.NewStyle().Background(backgroundColor).Width(lipgloss.Width(renderedCsPerMinute))

		renderedCs := cs.Render(csString)

		renderedCsSection := lipgloss.JoinVertical(lipgloss.Left, renderedCs, renderedCsPerMinute)

		renderedBody := lipgloss.JoinHorizontal(lipgloss.Top, renderedChampionSection, renderedKdaSection, renderedCsSection)

		container := lipgloss.NewStyle().Background(backgroundColor).Border(lipgloss.BlockBorder(), false, false, false, true).BorderForeground(borderColor).Padding(2).MarginBottom(1).Width(80)

		fmt.Println(container.Render(lipgloss.NewStyle().Background(backgroundColor).Render(renderedBody)))
	}
}
