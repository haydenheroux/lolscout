package model

import (
	"gorm.io/gorm"
	"time"
)

type Team struct {
	ID        string         `gorm:"primaryKey;column:id"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"`
	Name      string
	Players   []Player `gorm:"foreignKey:team_id"`
}

type Player struct {
	PUUID         string         `gorm:"primaryKey;column:puuid"`
	GameName      string         `gorm:"column:game_name"`
	TagLine       string         `gorm:"column:tag_line"`
	TeamID        *string        `gorm:"column:team_id"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index;column:deleted_at"`
	PlayerMetrics []MatchMetrics `gorm:"foreignKey:puuid"`
}

func (p *Player) AppendMatchMetrics(matchMetrics []*MatchMetrics) int {
	count := 0

	for _, metrics := range matchMetrics {
		skip := false

		for _, existing := range p.PlayerMetrics {
			if metrics.MatchID == existing.MatchID {
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		p.PlayerMetrics = append(p.PlayerMetrics, *metrics)
		count += 1
	}

	return count
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type MatchMetrics struct {
	gorm.Model

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
	RoleTop
	RoleJungle
	RoleMiddle
	RoleBottom
	RoleSupport
)
