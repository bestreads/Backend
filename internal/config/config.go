package config

import (
	"log"

	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	ApiProtocol                string `mapstructure:"API_PROTOCOL"`
	ApiDomain                  string `mapstructure:"API_DOMAIN"`
	ApiPort                    string `mapstructure:"API_PORT"`
	ApiProdPort                string `mapstructure:"API_PROD_PORT"`
	ApiBasePath                string `mapstructure:"API_BASE_PATH"`
	DebugLevel                 string `mapstructure:"DEBUG_LEVEL"`
	DBHost                     string `mapstructure:"DB_HOST"`
	DBPort                     string `mapstructure:"DB_PORT"`
	DBSslMode                  bool   `mapstructure:"DB_SSL_MODE"`
	DBName                     string `mapstructure:"DB_NAME"`
	DBUsername                 string `mapstructure:"DB_USERNAME"`
	DBPassword                 string `mapstructure:"DB_PASSWORD"`
	DataPath                   string `mapstructure:"DATA_PATH"`
	DevMode                    bool   `mapstructure:"DEV_MODE"`
	OpenRouterApiKey           string `mapstructure:"OPEN_ROUTER_API_KEY"`
	RefreshTokenSecretKey      string `mapstructure:"REFRESH_TOKEN_SECRET_KEY"`
	AccessTokenSecretKey       string `mapstructure:"ACCESS_TOKEN_SECRET_KEY"`
	AccessTokenDurationMinutes uint   `mapstructure:"ACCESS_TOKEN_DURATION_MINUTES"`
	RefreshTokenDurationDays   uint   `mapstructure:"REFRESH_TOKEN_DURATION_DAYS"`
	TokenRefreshPath           string `mapstructure:"TOKEN_REFRESH_PATH"`
	TokenSecureFlag            bool   `mapstructure:"TOKEN_SECURE_FLAG"`
	PaginationSteps            int    `mapstructure:"PAGINATION_STEPS"`
}

func Load() *Config {
	// Defaults
	viper.SetDefault("API_PROTOCOL", "http")
	viper.SetDefault("API_DOMAIN", "localhost")
	viper.SetDefault("API_PORT", "3000")
	viper.SetDefault("API_PROD_PORT", "443")
	viper.SetDefault("API_BASE_PATH", "/api")
	viper.SetDefault("DEBUG_LEVEL", "info")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSL_MODE", false)
	viper.SetDefault("DATA_PATH", "./store")
	viper.SetDefault("ACCESS_TOKEN_DURATION_MINUTES", 15)
	viper.SetDefault("REFRESH_TOKEN_DURATION_DAYS", 100)
	viper.SetDefault("TOKEN_REFRESH_PATH", "/api/v1/auth/refresh")
	viper.SetDefault("TOKEN_SECURE_FLAG", true)
	viper.SetDefault("DEV_MODE", false)
	viper.SetDefault("PAGINATION_STEPS", 25)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Read optional config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No config file found, falling back to ENV")
	}

	// ENV overrides file
	// Use for loop instead of AutomaticEnv() because of case sensitivity and upper case env vars
	v := reflect.TypeOf(Config{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag != "" {
			if err := viper.BindEnv(tag); err != nil {
				log.Printf("Failed to bind ENV for %s: %v", tag, err)
			}
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("Config Error: ", err)
	}

	return &cfg
}
