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

	PUUID string `gorm:"column:puuid;uniqueIndex:compositeIndex;"`

	MatchID string `gorm:"column:match_id;uniqueIndex:compositeIndex;"`

	StartTime time.Time

	Assists              int
	CS                   int
	CSPerMinute          float64
	ChampionName         Champion // TODO Rename to "Champion"
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

type Champion string

func (c Champion) String() string {
	return string(c)
}

type MatchType int

const (
	MatchTypeSummonersRift MatchType = iota
)

func (mt MatchType) String() string {
	switch mt {
	case MatchTypeSummonersRift:
		return "Summoner's Rift"
	default:
		return ""
	}
}

type Position int

const (
	Unknown Position = iota
	PositionTop
	PositionJungle
	PositionMiddle
	PositionBottom
	PositionSupport
)

var Positions = []Position{PositionTop, PositionJungle, PositionMiddle, PositionBottom, PositionSupport}

func (p Position) String() string {
	switch p {
	case Unknown:
		return "Unknown"
	case PositionTop:
		return "Top"
	case PositionJungle:
		return "Jungle"
	case PositionMiddle:
		return "Middle"
	case PositionBottom:
		return "Bottom"
	case PositionSupport:
		return "Support"
	default:
		return ""
	}
}
