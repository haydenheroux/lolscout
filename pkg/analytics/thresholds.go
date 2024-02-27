package analytics

import (
	"fmt"

	"github.com/montanaflynn/stats"
)

type Thresholds struct {
	Assists              float64
	CSPerMinute          float64
	ControlWardsPlaced   float64
	DamageDealtPerMinute float64
	DamageDealtShare     float64
	Deaths               float64
	KillParticipation    float64
	Kills                float64
	TurretsTaken         float64
	WardsKilled          float64
	WardsPlaced          float64
	WinRate              float64
}

func (t Thresholds) String() string {
	var s string

	s += fmt.Sprintln("Assists:", t.Assists)
	s += fmt.Sprintln("CSPerMinute:", t.CSPerMinute)
	s += fmt.Sprintln("ControlWardsPlaced:", t.ControlWardsPlaced)
	s += fmt.Sprintln("DamageDealtPerMinute:", t.DamageDealtPerMinute)
	s += fmt.Sprintln("DamageDealtShare:", t.DamageDealtShare)
	s += fmt.Sprintln("Deaths:", t.Deaths)
	s += fmt.Sprintln("KillParticipation:", t.KillParticipation)
	s += fmt.Sprintln("Kills:", t.Kills)
	s += fmt.Sprintln("TurretsTaken:", t.TurretsTaken)
	s += fmt.Sprintln("WardsKilled:", t.WardsKilled)
	s += fmt.Sprintln("WardsPlaced:", t.WardsPlaced)
	s += fmt.Sprintf("WinRate: %.4f\n", t.WinRate)

	return s
}

func GeneralThresholds() *Thresholds {
	return &Thresholds{
		Assists:              8.0,
		CSPerMinute:          7.0,
		ControlWardsPlaced:   1.0,
		DamageDealtPerMinute: 4500.0,
		DamageDealtShare:     0.25,
		Deaths:               5.0,
		KillParticipation:    0.5,
		Kills:                5.0,
		TurretsTaken:         2.0,
		WardsKilled:          2.0,
		WardsPlaced:          12.0,
		WinRate:              0.5,
	}
}

func PercentileTresholds(analytics Analytics, percentile float64) *Thresholds {
	percentileOf := func(n Norm) float64 {
		return stats.NormPpf(percentile, n.Mean, n.StdDev)
	}

	return &Thresholds{
		Assists:              percentileOf(analytics.Assists),
		CSPerMinute:          percentileOf(analytics.CSPerMinute),
		ControlWardsPlaced:   percentileOf(analytics.ControlWardsPlaced),
		DamageDealtPerMinute: percentileOf(analytics.DamageDealtPerMinute),
		DamageDealtShare:     percentileOf(analytics.DamageDealtShare),
		Deaths:               percentileOf(analytics.Deaths),
		KillParticipation:    percentileOf(analytics.KillParticipation),
		Kills:                percentileOf(analytics.Kills),
		TurretsTaken:         percentileOf(analytics.TurretsTaken),
		WardsKilled:          percentileOf(analytics.WardsKilled),
		WardsPlaced:          percentileOf(analytics.WardsPlaced),
		WinRate:              analytics.WinRate,
	}
}
