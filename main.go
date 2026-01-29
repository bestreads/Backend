package main

import (
	"os"

	"github.com/bestreads/Backend/internal/config"
	"github.com/bestreads/Backend/internal/server"
	"github.com/rs/zerolog"
)

// @title						Best Reads API
// @version						1.0
// @description					API for searching and managing books
// @termsOfService				https://github.com/bestreads/Backend/blob/main/LICENSE

// @contact.name				Best Reads Team
// @contact.url					https://github.com/bestreads/Backend/issues

// @license.name				The Unlicense
// @license.url					https://github.com/bestreads/Backend/blob/main/LICENSE

// @BasePath					/api
// @schemes						http https

// @securityDefinitions.apikey	BearerAuth
// @name						Authorization
// @in							header
// @description					Bearer token-based authentication. Use "Bearer {your-token}"

var (
	cfg    *config.Config
	logger zerolog.Logger
)

func init() {
	cfg = config.Load()

	// Configure zerolog logger
	logger = zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	level, errParsingDebugLevel := zerolog.ParseLevel(cfg.DebugLevel)
	if errParsingDebugLevel != nil {
		panic(errParsingDebugLevel)
	}
	zerolog.SetGlobalLevel(level)
}

func main() {
	server.Start(cfg, logger)
}
