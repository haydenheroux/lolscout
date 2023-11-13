package tui

import (
	"fmt"
	"lolscout/stats"

	"github.com/charmbracelet/lipgloss"
)

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

func RenderMatchParticipant(matchParticipant MatchParticipant) string {
	var backgroundColor, borderColor lipgloss.Color

	if matchParticipant.Win {
		backgroundColor = blueBackgroundColor
		borderColor = blueBorderColor
	} else {
		backgroundColor = redBackgroundColor
		borderColor = redBorderColor
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

	goldKda := whiteKda.Copy().Foreground(goldColor)
	blueKda := goldKda.Copy().Foreground(blueColor)
	greenKda := goldKda.Copy().Foreground(greenColor)

	kdaTextStyle := whiteKda

	if kdaValue >= stats.GOLD_KDA {
		kdaTextStyle = goldKda
	} else if kdaValue >= stats.BLUE_KDA {
		kdaTextStyle = blueKda
	} else if kdaValue >= stats.GREEN_KDA {
		kdaTextStyle = greenKda
	}

	kdaTextStyle = kdaTextStyle.PaddingRight(1)

	kdaText := fmt.Sprintf("%.2f:1 KDA", kdaValue)

	renderedKdaText := kdaTextStyle.Render(kdaText)

	killParticipationString := fmt.Sprintf("(%.0f%% KP)", matchParticipant.KillParticipation*100)

	killParticipation := lipgloss.NewStyle().Background(backgroundColor).PaddingRight(1)

	if matchParticipant.KillParticipation >= stats.KILL_PARTICIPATION {
		killParticipation = killParticipation.Foreground(redColor)
	}

	renderedKillParticipation := killParticipation.Render(killParticipationString)

	renderedKdaBottomSection := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedKdaText, renderedKillParticipation)

	kdaString := fmt.Sprintf("%d/%d/%d", matchParticipant.Kills, matchParticipant.Deaths, matchParticipant.Assists)

	kda := lipgloss.NewStyle().Background(backgroundColor).Bold(true).Width(lipgloss.Width(renderedKdaBottomSection))

	renderedKda := kda.Render(kdaString)

	// TODO Center renderedKda above bottom section
	renderedKdaSection := lipgloss.JoinVertical(lipgloss.Left, renderedKda, renderedKdaBottomSection)

	csString := fmt.Sprintf("%d CS", matchParticipant.CS)

	csPerMinuteString := fmt.Sprintf("%.1f CS/M", matchParticipant.CSPerMinute)

	csPerMinute := lipgloss.NewStyle().Background(backgroundColor)

	if matchParticipant.CSPerMinute >= stats.CS_PER_MINUTE {
		csPerMinute = csPerMinute.Foreground(redColor)
	}

	renderedCsPerMinute := csPerMinute.Render(csPerMinuteString)

	cs := lipgloss.NewStyle().Background(backgroundColor).Bold(true).Width(lipgloss.Width(renderedCsPerMinute))

	renderedCs := cs.Render(csString)

	renderedCsSection := lipgloss.JoinVertical(lipgloss.Left, renderedCs, renderedCsPerMinute)

	renderedBody := lipgloss.JoinHorizontal(lipgloss.Top, renderedChampionSection, renderedKdaSection, renderedCsSection)

	container := lipgloss.NewStyle().Background(backgroundColor).Border(lipgloss.BlockBorder(), false, false, false, true).BorderForeground(borderColor).Padding(2).MarginBottom(1).Width(80)

	return container.Render(lipgloss.NewStyle().Background(backgroundColor).Render(renderedBody))
}
