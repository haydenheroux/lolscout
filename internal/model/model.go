package model

import (
	"gorm.io/gorm"
	"time"
)

type Team struct {
	gorm.Model
	Name    string
	Players []Player
}

type Player struct {
	PUUID         string `gorm:"primaryKey;column:puuid"`
	GameName      string `gorm:"column:game_name"`
	TagLine       string `gorm:"column:tag_line"`
	TeamID        *uint
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index;column:deleted_at"`
	PlayerMetrics []MatchMetrics `gorm:"foreignKey:puuid"`
}

type MatchMetrics struct {
	gorm.Model

	// TODO Determine best method for handling duplicate metrics
	PUUID   string `gorm:"column:puuid;uniqueIndex:compositeIndex;"`
	MatchID string `gorm:"column:match_id;uniqueIndex:compositeIndex;"`

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
