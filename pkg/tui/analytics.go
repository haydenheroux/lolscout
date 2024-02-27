package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/haydenheroux/lolscout/pkg/analytics"
)

type Analytics struct {
	Analytics  *analytics.Analytics
	Thresholds *analytics.Thresholds
}

type analyticsTheme struct {
	gap int
}

func (a Analytics) theme() analyticsTheme {
	return analyticsTheme{
		gap: 1,
	}
}

func (a Analytics) View() string {
	var sb strings.Builder

	sb.WriteString(renderNorm("Assists", a.Analytics.Assists, exceeds(a.Thresholds.Assists, redColor)))
	sb.WriteString(renderNorm("CSPerMinute", a.Analytics.CSPerMinute, exceeds(a.Thresholds.CSPerMinute, redColor)))
	sb.WriteString(renderNorm("ControlWardsPlaced", a.Analytics.ControlWardsPlaced, exceeds(a.Thresholds.ControlWardsPlaced, redColor)))
	sb.WriteString(renderNorm("DamageDealtPerMinute", a.Analytics.DamageDealtPerMinute, exceeds(a.Thresholds.DamageDealtPerMinute, redColor)))
	sb.WriteString(renderPercent("DamageDealtShare", a.Analytics.DamageDealtShare.Mean, exceeds(a.Thresholds.DamageDealtShare, redColor)))
	sb.WriteString(renderNorm("Deaths", a.Analytics.Deaths, exceeds(a.Thresholds.Deaths, redColor)))
	sb.WriteString(renderPercent("KillParticipation", a.Analytics.KillParticipation.Mean, exceeds(a.Thresholds.KillParticipation, redColor)))
	sb.WriteString(renderNorm("Kills", a.Analytics.Kills, exceeds(a.Thresholds.Kills, redColor)))
	sb.WriteString(renderNorm("TurretsTaken", a.Analytics.TurretsTaken, exceeds(a.Thresholds.TurretsTaken, redColor)))
	sb.WriteString(renderNorm("WardsKilled", a.Analytics.WardsKilled, exceeds(a.Thresholds.WardsKilled, redColor)))
	sb.WriteString(renderNorm("WardsPlaced", a.Analytics.WardsPlaced, exceeds(a.Thresholds.WardsPlaced, redColor)))
	sb.WriteString(renderPercent("WinRate", a.Analytics.WinRate, exceeds(a.Thresholds.WinRate, redColor)))

	return strings.TrimSpace(sb.String())
}

type colorer func(value float64) lipgloss.Color

func exceeds(threshold float64, color lipgloss.Color) colorer {
	defaultColor := whiteColor

	return func(value float64) lipgloss.Color {
		if value > threshold {
			return color
		}

		return defaultColor
	}
}

func renderNorm(key string, value analytics.Norm, colorer colorer) string {
	meanStyle := lipgloss.NewStyle().Foreground(colorer(value.Mean))

	meanStr := fmt.Sprintf("%.2f", value.Mean)

	renderedMeanStr := meanStyle.Render(meanStr)

	keyStyle := lipgloss.NewStyle().Foreground(whiteColor)

	renderedKeyStr := keyStyle.Render(key)

	return fmt.Sprintf("%s: %s\n", renderedKeyStr, renderedMeanStr)
}

func renderPercent(key string, percent float64, colorer colorer) string {
	percentStyle := lipgloss.NewStyle().Foreground(colorer(percent))

	percentStr := fmt.Sprintf("%.2f%%", percent*100)

	renderedPercentStr := percentStyle.Render(percentStr)

	keyStyle := lipgloss.NewStyle().Foreground(whiteColor)

	renderedKeyStr := keyStyle.Render(key)

	return fmt.Sprintf("%s: %s\n", renderedKeyStr, renderedPercentStr)
}
