package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ApiPort     string `mapstructure:"API_PORT"`
	ApiBasePath string `mapstructure:"API_BASE_PATH"`
	DebugLevel  string `mapstructure:"DEBUG_LEVEL"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      string `mapstructure:"DB_PORT"`
	DBSslMode   bool   `mapstructure:"DB_SSL_MODE"`
	DBName      string `mapstructure:"DB_NAME"`
	DBUsername  string `mapstructure:"DB_USERNAME"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
}

func Load() *Config {
	// Defaults
	viper.SetDefault("API_PORT", "3000")
	viper.SetDefault("API_BASE_PATH", "/api")
	viper.SetDefault("DEBUG_LEVEL", "info")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSL_MODE", false)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Read optional config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No config file found, falling back to ENV")
	}

	// ENV overrides file
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("Config Error: ", err)
	}

	return &cfg
}
