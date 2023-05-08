package db

import (
	"contact-go/config"
	"contact-go/helper/apperrors"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormDatabase(cfg *config.Config) (*gorm.DB, error) {
	if cfg.Database.URL == "" {
		return nil, apperrors.NewAppError(apperrors.ErrDbUrlNotExist)
	}

	var gormLogger logger.Interface
	gormLogger = logger.Default.LogMode(logger.Silent)
	if cfg.Debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	gormConf := new(gorm.Config)
	gormConf.Logger = gormLogger
	gormConf.PrepareStmt = true
	gormConf.SkipDefaultTransaction = true

	db, err := gorm.Open(postgres.Open(cfg.Database.URL), gormConf)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
