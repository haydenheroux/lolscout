package analytics

import (
	"fmt"

	"github.com/haydenheroux/lolscout/pkg/model"
	"github.com/montanaflynn/stats"
)

type Norm struct {
	Mean   float64
	StdDev float64
}

func (n Norm) String() string {
	return fmt.Sprintf("N(μ: %.4f, σ: %.4f)", n.Mean, n.StdDev)
}

func calculateNorm(xs interface{}) Norm {
	data := stats.LoadRawData(xs)

	mean, _ := data.Mean()
	stdDev, _ := data.StandardDeviation()

	return Norm{
		Mean:   mean,
		StdDev: stdDev,
	}
}

type Analytics struct {
	Assists              Norm
	CSPerMinute          Norm
	ControlWardsPlaced   Norm
	DamageDealtPerMinute Norm
	DamageDealtShare     Norm
	Deaths               Norm
	KillParticipation    Norm
	Kills                Norm
	Size                 int
	TurretsTaken         Norm
	WardsKilled          Norm
	WardsPlaced          Norm
	WinRate              float64
}

func (a Analytics) String() string {
	var s string

	s += fmt.Sprintln("Assists:", a.Assists)
	s += fmt.Sprintln("CSPerMinute:", a.CSPerMinute)
	s += fmt.Sprintln("ControlWardsPlaced:", a.ControlWardsPlaced)
	s += fmt.Sprintln("DamageDealtPerMinute:", a.DamageDealtPerMinute)
	s += fmt.Sprintln("DamageDealtShare:", a.DamageDealtShare)
	s += fmt.Sprintln("Deaths:", a.Deaths)
	s += fmt.Sprintln("KillParticipation:", a.KillParticipation)
	s += fmt.Sprintln("Kills:", a.Kills)
	s += fmt.Sprintln("Size:", a.Size)
	s += fmt.Sprintln("TurretsTaken:", a.TurretsTaken)
	s += fmt.Sprintln("WardsKilled:", a.WardsKilled)
	s += fmt.Sprintln("WardsPlaced:", a.WardsPlaced)
	s += fmt.Sprintf("WinRate: %.4f\n", a.WinRate)

	return s
}

func Analyze(metrics []model.MatchMetrics) *Analytics {
	assists := make([]int, len(metrics))
	csPerMinute := make([]float64, len(metrics))
	controlWardsPlaced := make([]int, len(metrics))
	damageDealtPerMinute := make([]float64, len(metrics))
	damageDealtShare := make([]float64, len(metrics))
	deaths := make([]int, len(metrics))
	killParticipation := make([]float64, len(metrics))
	kills := make([]int, len(metrics))
	turretsTaken := make([]int, len(metrics))
	wardsKilled := make([]int, len(metrics))
	wardsPlaced := make([]int, len(metrics))
	wins := make([]bool, len(metrics))

	for i, metric := range metrics {
		assists[i] = metric.Assists
		csPerMinute[i] = metric.CSPerMinute
		controlWardsPlaced[i] = metric.ControlWardsPlaced
		damageDealtPerMinute[i] = metric.DamageDealtPerMinute
		damageDealtShare[i] = metric.DamageDealtShare
		deaths[i] = metric.Deaths
		killParticipation[i] = metric.KillParticipation
		kills[i] = metric.Kills
		turretsTaken[i] = metric.TurretsTaken
		wardsKilled[i] = metric.WardsKilled
		wardsPlaced[i] = metric.WardsPlaced
		wins[i] = metric.Win
	}

	assistsNorm := calculateNorm(assists)
	csPerMinuteNorm := calculateNorm(csPerMinute)
	controlWardsPlacedNorm := calculateNorm(controlWardsPlaced)
	damageDealtPerMinuteNorm := calculateNorm(damageDealtPerMinute)
	damageDealtShareNorm := calculateNorm(damageDealtShare)
	deathsNorm := calculateNorm(deaths)
	killParticipationNorm := calculateNorm(killParticipation)
	killsNorm := calculateNorm(kills)
	turretsTakenNorm := calculateNorm(turretsTaken)
	wardsKilledNorm := calculateNorm(wardsKilled)
	wardsPlacedNorm := calculateNorm(wardsPlaced)
	winRate := percentTrue(wins)

	return &Analytics{
		Assists:              assistsNorm,
		CSPerMinute:          csPerMinuteNorm,
		ControlWardsPlaced:   controlWardsPlacedNorm,
		DamageDealtPerMinute: damageDealtPerMinuteNorm,
		DamageDealtShare:     damageDealtShareNorm,
		Deaths:               deathsNorm,
		KillParticipation:    killParticipationNorm,
		Kills:                killsNorm,
		Size:                 len(metrics),
		TurretsTaken:         turretsTakenNorm,
		WardsKilled:          wardsKilledNorm,
		WardsPlaced:          wardsPlacedNorm,
		WinRate:              winRate,
	}
}

func percentTrue(slice []bool) float64 {
	trues := 0

	for _, value := range slice {
		if value {
			trues++
		}
	}

	return float64(trues) / float64(len(slice))
}

type AnalyticsByChampion map[model.Champion]*Analytics

func AnalyzeByChampion(metrics []model.MatchMetrics) AnalyticsByChampion {
	metricsByChampion := byChampion(metrics)

	analyticsByChampion := make(AnalyticsByChampion)

	for champion, metrics := range metricsByChampion {
		analyticsByChampion[champion] = Analyze(metrics)
	}

	return analyticsByChampion
}

type championMetrics map[model.Champion][]model.MatchMetrics

func byChampion(metrics []model.MatchMetrics) championMetrics {
	championMetrics := make(championMetrics)

	for _, metric := range metrics {
		champion := metric.Champion

		if _, ok := championMetrics[champion]; !ok {
			championMetrics[champion] = []model.MatchMetrics{metric}
		} else {
			championMetrics[champion] = append(championMetrics[champion], metric)
		}
	}

	return championMetrics
}

type AnalyticsByPosition map[model.Position]*Analytics

func AnalyzeByPosition(metrics []model.MatchMetrics) AnalyticsByPosition {
	metricsByPosition := byPosition(metrics)

	analyticsByPosition := make(AnalyticsByPosition)

	for position, metrics := range metricsByPosition {
		analyticsByPosition[position] = Analyze(metrics)
	}

	return analyticsByPosition
}

type positionMetrics map[model.Position][]model.MatchMetrics

func byPosition(metrics []model.MatchMetrics) positionMetrics {
	positionMetrics := make(positionMetrics)

	for _, metric := range metrics {
		position := metric.Position

		if _, ok := positionMetrics[position]; !ok {
			positionMetrics[position] = []model.MatchMetrics{metric}
		} else {
			positionMetrics[position] = append(positionMetrics[position], metric)
		}
	}

	return positionMetrics
}
