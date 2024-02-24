package model

import "gorm.io/gorm"

type MatchMetrics struct {
	gorm.Model
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
	MatchType            MatchType
	Position             Position
	TurretsTaken         int
	WardsKilled          int
	WardsPlaced          int
	Win                  bool
	PlayerID             uint
}

type MatchType int

const (
	SummonersRift MatchType = iota
)

type Position int

const (
	Unknown Position = iota
	Top
	Jungle
	Middle
	Bottom
	Support
)
