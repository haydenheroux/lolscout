package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/haydenheroux/lolscout/internal/model"
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

func (dbc client) CreateOrUpdatePlayer(player *model.Player) error {
	return dbc.DB.Save(player).Error
}
