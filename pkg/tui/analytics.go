package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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

func createTable() *table.Table {
	t := table.New().Border(lipgloss.NormalBorder()).BorderStyle(lipgloss.NewStyle()).BorderRow(true).BorderColumn(true)

	t.StyleFunc(func(row, col int) lipgloss.Style {
		switch {
		case row == 0:
			return lipgloss.NewStyle().Bold(true).Foreground(draculaForegroundWhite)
		default:
			return lipgloss.NewStyle().Foreground(draculaForegroundWhite)
		}
	})

	t.Width(50)

	return t
}

func (a Analytics) View(title string) string {
	t := createTable()

	t.Headers(title)

	t.Row("Assists", fmt.Sprintf("%.2f", a.Analytics.Assists.Mean))
	t.Row("CSPerMinute", fmt.Sprintf("%.2f", a.Analytics.CSPerMinute.Mean))
	t.Row("ControlWardsPlaced", fmt.Sprintf("%.2f", a.Analytics.ControlWardsPlaced.Mean))
	t.Row("DamageDealtPerMinute", fmt.Sprintf("%.2f", a.Analytics.DamageDealtPerMinute.Mean))
	t.Row("DamageDealtShare", fmt.Sprintf("%.2f", a.Analytics.DamageDealtShare.Mean))
	t.Row("Deaths", fmt.Sprintf("%.2f", a.Analytics.Deaths.Mean))
	t.Row("KillParticipation", fmt.Sprintf("%.2f", a.Analytics.KillParticipation.Mean))
	t.Row("Kills", fmt.Sprintf("%.2f", a.Analytics.Kills.Mean))
	t.Row("TurretsTaken", fmt.Sprintf("%.2f", a.Analytics.TurretsTaken.Mean))
	t.Row("WardsKilled", fmt.Sprintf("%.2f", a.Analytics.WardsKilled.Mean))
	t.Row("WardsPlaced", fmt.Sprintf("%.2f", a.Analytics.WardsPlaced.Mean))
	t.Row("WinRate", fmt.Sprintf("%.2f", a.Analytics.WinRate))

	return t.String()
}

func getMany(get func(a *analytics.Analytics) float64, as []*analytics.Analytics) []float64 {
	var result []float64

	for _, a := range as {
		result = append(result, get(a))
	}

	return result
}

func format(ns []float64) []string {
	var result []string

	for _, n := range ns {
		result = append(result, fmt.Sprintf("%.2f", n))
	}

	return result
}

func ViewAnalytics(headers []string, t *analytics.Thresholds, as []*analytics.Analytics) string {
	table := createTable()

	table.Headers(headers...)

	getAssists := func(a *analytics.Analytics) float64 {
		return a.Assists.Mean
	}

	assists := []string{"Assists"}

	assists = append(assists, format([]float64{t.Assists})...)
	assists = append(assists, format(getMany(getAssists, as))...)

	table.Row(assists...)

	return table.String()
}

type colorer func(value float64) lipgloss.Color

func exceeds(threshold float64, color lipgloss.Color) colorer {
	defaultColor := opggWhite

	return func(value float64) lipgloss.Color {
		if value > threshold {
			return color
		}

		return defaultColor
	}
}
