package tui

import (
	"fmt"
	"lolscout/stats"

	"github.com/charmbracelet/lipgloss"
)

type MatchParticipantModel struct {
	ChampionName      string
	Level             int
	Kills             int
	Deaths            int
	Assists           int
	KillParticipation float64
	CS                int
	CSPerMinute       float64
	Win               bool
	MatchType         string
	DurationMinutes   int
}

type theme struct {
	background lipgloss.Color
	border     lipgloss.Color
	gap        int
}

func (mp MatchParticipantModel) theme() theme {
	if mp.Win {
		return theme{
			background: blueBackgroundColor,
			border:     blueBorderColor,
			gap:        1,
		}
	} else {
		return theme{
			background: redBackgroundColor,
			border:     redBorderColor,
			gap:        1,
		}
	}
}

func (mp MatchParticipantModel) infoSectionView() string {
	background := mp.theme().background

	durationString := fmt.Sprintf("%dm", mp.DurationMinutes)

	width := max(len(mp.MatchType), len(durationString))

	matchType := lipgloss.NewStyle().Background(background).Bold(true).Width(width)
	renderedMatchType := matchType.Render(mp.MatchType)

	duration := lipgloss.NewStyle().Background(background).Width(width)
	renderedDuration := duration.Render(durationString)

	return lipgloss.JoinVertical(lipgloss.Left, renderedMatchType, renderedDuration)
}

func (mp MatchParticipantModel) championView() string {
	background := mp.theme().background

	levelString := fmt.Sprintf("Lvl. %d", mp.Level)

	width := max(len(levelString), len(mp.ChampionName))

	championName := lipgloss.NewStyle().Background(background).Bold(true).Width(width)

	renderedChampionName := championName.Render(mp.ChampionName)

	level := lipgloss.NewStyle().Background(background).Width(width)

	renderedLevel := level.Render(levelString)

	return lipgloss.JoinVertical(lipgloss.Left, renderedChampionName, renderedLevel)
}

func (mp MatchParticipantModel) kdaView() string {
	background := mp.theme().background

	kdRatio := float64(mp.Kills + mp.Assists)

	if mp.Deaths > 0 {
		kdRatio /= float64(mp.Deaths)
	}

	whiteKda := lipgloss.NewStyle().Background(background)

	goldKda := whiteKda.Copy().Foreground(goldColor)
	blueKda := goldKda.Copy().Foreground(blueColor)
	greenKda := goldKda.Copy().Foreground(greenColor)

	kdaTextStyle := whiteKda

	if kdRatio >= stats.GOLD_KDA {
		kdaTextStyle = goldKda
	} else if kdRatio >= stats.BLUE_KDA {
		kdaTextStyle = blueKda
	} else if kdRatio >= stats.GREEN_KDA {
		kdaTextStyle = greenKda
	}

	kdaTextStyle = kdaTextStyle.PaddingRight(1)

	kdaText := fmt.Sprintf("%.2f:1 KDA", kdRatio)

	renderedKdaText := kdaTextStyle.Render(kdaText)

	killParticipationString := fmt.Sprintf("(%.0f%% KP)", mp.KillParticipation*100)

	killParticipation := lipgloss.NewStyle().Background(background)

	if mp.KillParticipation >= stats.KILL_PARTICIPATION {
		killParticipation = killParticipation.Foreground(redColor)
	}

	renderedKillParticipation := killParticipation.Render(killParticipationString)

	// TODO Take width into account
	renderedKdaBottomSection := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedKdaText, renderedKillParticipation)

	kdaString := fmt.Sprintf("%d/%d/%d", mp.Kills, mp.Deaths, mp.Assists)

	width := max(len(kdaString), lipgloss.Width(renderedKdaBottomSection))

	kda := lipgloss.NewStyle().Background(background).Bold(true).Width(width) //.Align(lipgloss.Center)

	renderedKda := kda.Render(kdaString)

	return lipgloss.JoinVertical(lipgloss.Left, renderedKda, renderedKdaBottomSection)
}

func (mp MatchParticipantModel) creepScoreView() string {
	background := mp.theme().background

	csString := fmt.Sprintf("%d CS", mp.CS)

	csPerMinuteString := fmt.Sprintf("%.1f CS/M", mp.CSPerMinute)

	csPerMinute := lipgloss.NewStyle().Background(background)

	if mp.CSPerMinute >= stats.CS_PER_MINUTE {
		csPerMinute = csPerMinute.Foreground(redColor)
	}

	renderedCsPerMinute := csPerMinute.Render(csPerMinuteString)

	cs := lipgloss.NewStyle().Background(background).Bold(true).Width(lipgloss.Width(renderedCsPerMinute))

	renderedCs := cs.Render(csString)

	return lipgloss.JoinVertical(lipgloss.Left, renderedCs, renderedCsPerMinute)
}

func (mp MatchParticipantModel) View() string {
	theme := mp.theme()
	background := theme.background
	border := theme.border

	renderedBody := theme.joinHorizontal(mp.infoSectionView(), mp.championView(), mp.kdaView(), mp.creepScoreView())

	container := lipgloss.NewStyle().Background(background).Border(lipgloss.BlockBorder(), false, false, false, true).BorderForeground(border).Padding(2).MarginBottom(1).Width(60)

	return container.Render(renderedBody)
}

func (theme theme) joinHorizontal(sections ...string) string {
	if len(sections) == 0 {
		return ""
	}

	style := lipgloss.NewStyle().Background(theme.background).PaddingRight(theme.gap)

	joined := make([]string, len(sections))

	for i, s := range sections {
		if i != len(sections)-1 {
			s = style.Render(s)
		}

		joined[i] = s
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, joined...)
}
