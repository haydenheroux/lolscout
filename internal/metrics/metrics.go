package metrics

import (
	"encoding/csv"
	"io"
	"strconv"
)

type MatchMetrics struct {
	Assists              int
	CS                   int
	CSPerMinute          float64
	ChampionName         string
	ControlWardsPlaced   int
	DamageDealt          int
	DamageDealtPerMinute float64
	DamageDealtShare     float64
	Deaths               int
	DurationMinutes      float64
	KillParticipation    float64
	Kills                int
	Level                int
	MatchType            string
	Position             string
	TurretsTaken         int
	WardsKilled          int
	WardsPlaced          int
	Win                  bool
}

type MetricsCollection []*MatchMetrics

func formatBool(b bool) string {
	return strconv.FormatBool(b)
}

func formatInt(i int) string {
	return strconv.Itoa(i)
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func (m MatchMetrics) header() []string {
	return []string{
		"assists",
		"cs",
		"cs/m",
		"champion",
		"controlWards",
		"dmg",
		"dmg/m",
		"dmg%",
		"deaths",
		"durationMinutes",
		"kp",
		"kills",
		"level",
		"matchType",
		"position",
		"turrets",
		"wardsKilled",
		"wardsPlaced",
		"win",
	}
}

func (m MatchMetrics) row() []string {
	return []string{
		formatInt(m.Assists),
		formatInt(m.CS),
		formatFloat(m.CSPerMinute),
		m.ChampionName,
		formatInt(m.ControlWardsPlaced),
		formatInt(m.DamageDealt),
		formatFloat(m.DamageDealtPerMinute),
		formatFloat(m.DamageDealtShare),
		formatInt(m.Deaths),
		formatFloat(m.DurationMinutes),
		formatFloat(m.KillParticipation),
		formatInt(m.Kills),
		formatInt(m.Level),
		m.MatchType,
		m.Position,
		formatInt(m.TurretsTaken),
		formatInt(m.WardsKilled),
		formatInt(m.WardsPlaced),
		formatBool(m.Win),
	}
}

func (mc MetricsCollection) CSV(w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	headerWritten := false

	for _, metrics := range mc {
		if !headerWritten {
			err := writer.Write(metrics.header())

			if err != nil {
				return err
			}

			headerWritten = true
		}

		err := writer.Write(metrics.row())
		if err != nil {
			return err
		}
	}

	return nil
}
