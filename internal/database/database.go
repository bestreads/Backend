package database

import (
	"context"
	"fmt"

	"github.com/bestreads/Backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// db muss nicht überall als argument übergeben werden
// es ist nicht funktional aber w/e
var (
	GlobalDB *gorm.DB
)

func SetupDatabase(cfg *config.Config, ctx context.Context) (*gorm.DB, error) {
	sslMode := "disable"
	if cfg.DBSslMode {
		sslMode = "require"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", cfg.DBHost, cfg.DBUsername, cfg.DBPassword, cfg.DBName, cfg.DBPort, sslMode)

	var err error

	GlobalDB, err = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)},
	)

	if err != nil {
		return nil, err
	}

	if err = GlobalDB.AutoMigrate(&User{}, &UserMeta{}); err != nil {
		return nil, err
	}

	if err = GlobalDB.AutoMigrate(&Book{}); err != nil {
		return nil, err
	}

	if err = GlobalDB.AutoMigrate(&RelBookUser{}); err != nil {
		return nil, err
	}

	if err = GlobalDB.AutoMigrate(&Post{}); err != nil {
		return nil, err
	}

	if err := insertDemoData(GlobalDB, ctx); err != nil {
		println("soft error mit den demodaten")
	}

	return GlobalDB, nil
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
	description string,
	releaseDate uint64) error {

	b := &Book{ISBN: safeIsbn, Title: title, Author: author, Description: description, ReleaseDate: releaseDate}
	return gorm.G[Book](db).Create(ctx, b)
}

func CreateUserBookRel(db *gorm.DB, ctx context.Context, uid uint, bid uint, s state) error {
	return gorm.G[RelBookUser](db).Create(ctx, &RelBookUser{UserID: uid, BookID: bid, State: s})
}

func validateISBN(unsafeIsbn string) (string, error) {
	return unsafeIsbn, nil
}

func insertDemoData(db *gorm.DB, ctx context.Context) error {

	for i := range 10 {
		if err := CreateUser(db, ctx, fmt.Sprintf("%d@test.com", i), "aaaaaaaaaaa"); err != nil {
			return err
		}
	}

	for i := range 10 {
		if err := CreateBook(db, ctx, fmt.Sprintf(
			"978-0-439-0234%d-2", i),
			fmt.Sprintf("test%d", i),
			fmt.Sprintf("test%d", i),
			fmt.Sprintf("test%d", i),
			uint64(i)); err != nil {
			return err
		}
	}

	for i := range 2 {
		if err := CreateUserBookRel(db, ctx, uint(i+1), uint(i+1), state(i)); err != nil {
			return err
		}
	}

	return nil
}
