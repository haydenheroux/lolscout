package db

import (
	"github.com/haydenheroux/lolscout/pkg/analytics"
	"github.com/haydenheroux/lolscout/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type client struct {
	DB *gorm.DB
}

func CreateClient(dsn string) (*client, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		return &client{}, err
	}

	err = db.AutoMigrate(&model.Team{}, &model.Player{}, &model.MatchMetrics{})

	if err != nil {
		return &client{}, err
	}

	return &client{DB: db}, nil
}

func (dbc client) CreateOrUpdateTeam(team *model.Team) error {
	return dbc.DB.Save(team).Error
}

func (dbc client) GetAllTeams() ([]*model.Team, error) {
	var teams []*model.Team
	if err := dbc.DB.Find(&teams).Error; err != nil {
		return nil, err
	}
	return teams, nil
}

func (dbc client) GetTeamByID(id string) (*model.Team, error) {
	var team model.Team
	if err := dbc.DB.Preload("Players").First(&team, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (dbc client) CreateOrUpdatePlayer(player *model.Player) error {
	return dbc.DB.Save(player).Error
}

func (dbc client) GetPlayerByPUUID(puuid string) (*model.Player, error) {
	var player model.Player
	if err := dbc.DB.Model(&model.Player{}).Preload("PlayerMetrics").First(&player, "puuid = ?", puuid).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func (dbc client) GetMatchIDsForPUUID(puuid string) ([]string, error) {
	var matchMetrics []model.MatchMetrics

	if err := dbc.DB.Model(&model.MatchMetrics{}).Select("match_id").Where("puuid = ?", puuid).Find(&matchMetrics).Error; err != nil {
		return nil, err
	}

	var matchIDs []string

	for _, metric := range matchMetrics {
		matchIDs = append(matchIDs, metric.MatchID)
	}

	return matchIDs, nil
}

func (dbc client) GetMetricsForPosition(position model.Position) ([]model.MatchMetrics, error) {
	var metrics []model.MatchMetrics

	if err := dbc.DB.Model(&model.MatchMetrics{}).Where("position = ?", position).Find(&metrics).Error; err != nil {
		return nil, err
	}

	return metrics, nil
}

func (dbc client) GetMetricsForChampion(champion model.Champion) ([]model.MatchMetrics, error) {
	var metrics []model.MatchMetrics

	if err := dbc.DB.Model(&model.MatchMetrics{}).Where("champion_name = ?", champion).Find(&metrics).Error; err != nil {
		return nil, err
	}

	return metrics, nil
}

func (dbc client) GetPositionThresholds(percentile float64) (map[model.Position]*analytics.Thresholds, error) {
	result := make(map[model.Position]*analytics.Thresholds)

	for _, position := range model.Positions {
		metrics, err := dbc.GetMetricsForPosition(position)
		if err != nil {
			return result, err
		}

		positionAnalytics := analytics.Analyze(metrics)

		result[position] = analytics.PercentileTresholds(*positionAnalytics, percentile)
	}

	return result, nil
}
