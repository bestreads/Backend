package main

import (
	"os"

	"github.com/bestreads/Backend/internal/config"
	"github.com/bestreads/Backend/internal/server"
	"github.com/rs/zerolog"
)

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
