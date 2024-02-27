package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/haydenheroux/lolscout/pkg/model"
)

type MatchMetrics struct {
	Metrics model.MatchMetrics
}

type matchMetricsTheme struct {
	background lipgloss.Color
	border     lipgloss.Color
	gap        int
}

func (m MatchMetrics) theme() matchMetricsTheme {
	if m.Metrics.Win {
		return matchMetricsTheme{
			background: blueBackgroundColor,
			border:     blueBorderColor,
			gap:        1,
		}
	} else {
		return matchMetricsTheme{
			background: redBackgroundColor,
			border:     redBorderColor,
			gap:        1,
		}
	}
}

func (m MatchMetrics) infoSectionView() string {
	background := m.theme().background

	durationString := fmt.Sprintf("%dm", int(m.Metrics.DurationMinutes))

	minWidth := 10
	width := max(len(m.Metrics.MatchType.String()), len(durationString), minWidth)

	matchType := lipgloss.NewStyle().Background(background).Bold(true).Width(width)
	renderedMatchType := matchType.Render(m.Metrics.MatchType.String())

	duration := lipgloss.NewStyle().Background(background).Width(width)
	renderedDuration := duration.Render(durationString)

	return lipgloss.JoinVertical(lipgloss.Left, renderedMatchType, renderedDuration)
}

func (m MatchMetrics) championView() string {
	background := m.theme().background

	levelString := fmt.Sprintf("Lvl. %d", m.Metrics.Level)

	minWidth := 10
	width := max(len(levelString), len(m.Metrics.ChampionName), minWidth)

	championName := lipgloss.NewStyle().Background(background).Bold(true).Width(width)

	renderedChampionName := championName.Render(m.Metrics.ChampionName.String())

	level := lipgloss.NewStyle().Background(background).Width(width)

	renderedLevel := level.Render(levelString)

	return lipgloss.JoinVertical(lipgloss.Left, renderedChampionName, renderedLevel)
}

func (m MatchMetrics) kdaView() string {
	background := m.theme().background

	kdRatio := float64(m.Metrics.Kills + m.Metrics.Assists)

	if m.Metrics.Deaths > 0 {
		kdRatio /= float64(m.Metrics.Deaths)
	}

	whiteKda := lipgloss.NewStyle().Background(background)

	// TODO
	// goldKda := whiteKda.Copy().Foreground(goldColor)
	// blueKda := goldKda.Copy().Foreground(blueColor)
	// greenKda := goldKda.Copy().Foreground(greenColor)

	kdaTextStyle := whiteKda

	// TODO
	// if kdRatio >= data.GOLD_KDA {
	// 	kdaTextStyle = goldKda
	// } else if kdRatio >= data.BLUE_KDA {
	// 	kdaTextStyle = blueKda
	// } else if kdRatio >= data.GREEN_KDA {
	// 	kdaTextStyle = greenKda
	// }

	kdaTextStyle = kdaTextStyle.PaddingRight(1)

	kdaText := fmt.Sprintf("%.2f:1 KDA", kdRatio)

	renderedKdaText := kdaTextStyle.Render(kdaText)

	killParticipationString := fmt.Sprintf("(%.0f%% KP)", m.Metrics.KillParticipation*100)

	killParticipation := lipgloss.NewStyle().Background(background)

	// TODO
	// if mp.KillParticipation >= data.KILL_PARTICIPATION {
	// 	killParticipation = killParticipation.Foreground(redColor)
	// }

	renderedKillParticipation := killParticipation.Render(killParticipationString)

	// TODO Take width into account
	renderedKdaBottomSection := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedKdaText, renderedKillParticipation)

	kdaString := fmt.Sprintf("%d/%d/%d", m.Metrics.Kills, m.Metrics.Deaths, m.Metrics.Assists)

	minWidth := 24
	width := max(len(kdaString), lipgloss.Width(renderedKdaBottomSection), minWidth)

	// TODO Sketchy
	kdaBottomSection := lipgloss.NewStyle().Background(background).Width(width)

	rerenderedKdaBottomSection := kdaBottomSection.Render(renderedKdaBottomSection)

	kda := lipgloss.NewStyle().Background(background).Bold(true).Width(width) //.Align(lipgloss.Center)

	renderedKda := kda.Render(kdaString)

	return lipgloss.JoinVertical(lipgloss.Left, renderedKda, rerenderedKdaBottomSection)
}

func (m MatchMetrics) creepScoreView() string {
	background := m.theme().background

	csString := fmt.Sprintf("%d CS", m.Metrics.CS)

	csPerMinuteString := fmt.Sprintf("%.1f/m", m.Metrics.CSPerMinute)

	csPerMinute := lipgloss.NewStyle().Background(background)

	// TODO
	// if mp.CSPerMinute >= data.CS_PER_MINUTE {
	// 	csPerMinute = csPerMinute.Foreground(redColor)
	// }

	renderedCsPerMinute := csPerMinute.Render(csPerMinuteString)

	minWidth := 10
	width := max(len(csString), lipgloss.Width(renderedCsPerMinute), minWidth)

	csPerMinute = csPerMinute.Width(width)

	rerenderedCsPerMinute := csPerMinute.Render(renderedCsPerMinute)

	cs := lipgloss.NewStyle().Background(background).Bold(true).Width(width)

	renderedCs := cs.Render(csString)

	return lipgloss.JoinVertical(lipgloss.Left, renderedCs, rerenderedCsPerMinute)
}

func (m MatchMetrics) View() string {
	theme := m.theme()
	background := theme.background
	border := theme.border

	renderedBody := theme.joinHorizontal(m.infoSectionView(), m.championView(), m.kdaView(), m.creepScoreView())

	container := lipgloss.NewStyle().Background(background).Border(lipgloss.BlockBorder(), false, false, false, true).BorderForeground(border).Padding(2).MarginBottom(1)

	return container.Render(renderedBody)
}

func (theme matchMetricsTheme) joinHorizontal(sections ...string) string {
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
