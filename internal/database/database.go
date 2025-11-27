package database

import (
	"fmt"

	"github.com/bestreads/Backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase(cfg *config.Config) (*gorm.DB, error) {
	sslMode := "disable"
	if cfg.DBSslMode {
		sslMode = "require"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.DBHost, cfg.DBUsername, cfg.DBPassword, cfg.DBName, cfg.DBPort, sslMode)

	// es fängt schon wieder an
	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)},
	)

	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&Book{}); err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&RelBookUser{}); err != nil {
		return nil, err
	}

	return db, nil
}
