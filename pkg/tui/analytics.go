package tui

import (
	"fmt"
	"reflect"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/haydenheroux/lolscout/pkg/analytics"
)

func createTable() *table.Table {
	t := table.New().Border(lipgloss.NormalBorder()).BorderStyle(lipgloss.NewStyle()).BorderRow(true).BorderColumn(true)

	t.StyleFunc(func(row, col int) lipgloss.Style {
		switch {
		case row == 0:
			return lipgloss.NewStyle().Bold(true).Foreground(draculaForegroundWhite).Align(lipgloss.Center)
		default:
			return lipgloss.NewStyle().Foreground(draculaForegroundWhite)
		}
	})

	t.Width(80)

	return t
}

func ViewAnalytics(title string, columns []string, analytics ...*analytics.AnalyticsSnapshot) string {
	t := createTable()

	t.Headers(append([]string{title}, columns...)...)

	s := reflect.ValueOf(analytics[0]).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		fieldName := typeOfT.Field(i).Name
		rowValues := make([]string, len(analytics)+1)
		rowValues[0] = fieldName

		for j, a := range analytics {
			s1 := reflect.ValueOf(a).Elem()

			rowValues[j+1] = fmt.Sprintf("%.2f", s1.Field(i).Float())
		}

		t.Row(rowValues...)
	}

	return t.String()
}
