package db

import (
	"github.com/haydenheroux/lolscout/internal/model"
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

func (dbc client) GetTeamByID(id uint) (*model.Team, error) {
	var team model.Team
	if err := dbc.DB.First(&team, id).Error; err != nil {
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
