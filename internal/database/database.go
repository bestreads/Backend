package database

import (
	"context"
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

func CreateUser(db *gorm.DB, ctx context.Context, mail string, hash string) error {
	return gorm.G[User](db).Create(ctx, &User{Email: mail, Password_hash: hash})
}

func CreateBook(
	db *gorm.DB,
	ctx context.Context,
	safeIsbn string,
	title string,
	author string,
	description string) error {

	b := &Book{ISBN: safeIsbn, Title: title, Author: author, Description: description}
	return gorm.G[Book](db).Create(ctx, b)
}

func CreateUserBookRel(db *gorm.DB, ctx context.Context, uid uint, bid uint, s state) error {
	return gorm.G[RelBookUser](db).Create(ctx, &RelBookUser{UserID: uid, BookID: bid, State: s})
}

func validateISBN(unsafeIsbn string) (string, error) {
	return unsafeIsbn, nil
}
