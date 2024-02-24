package model

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	Name    string
	Players []Player
}

type Player struct {
	gorm.Model
	PUUID        string
	GameName     string
	TagLine      string
	TeamID       *uint
	MatchMetrics []MatchMetrics
}
