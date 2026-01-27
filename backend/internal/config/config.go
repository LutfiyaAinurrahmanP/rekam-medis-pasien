package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig
	Database   DatabaseConfig
	JWT        JWTConfig
	Pagination PaginationConfig
}

type AppConfig struct {
	Env  string
	Name string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Timezone string
}

type JWTConfig struct {
	Secret      string
	ExpiredTime time.Duration
}

type PaginationConfig struct {
	DefaultPageSize int
	MaxPageSize     int
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading config file: %v. Using environment variables.", err)
	}

	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_NAME", "Sirekam Medis API")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_TIMEZONE", "Asia/Jakarta")
	viper.SetDefault("JWT_EXPIRED_TIME", "24h")
	viper.SetDefault("DEFAULT_PAGE_SIZE", 10)
	viper.SetDefault("MAX_PAGE_SIZE", 100)

	jwtExpired, err := time.ParseDuration(viper.GetString("JWT_EXPIRED_TIME"))
	if err != nil {
		jwtExpired = 24 * time.Hour
	}

	config := &Config{
		App: AppConfig{
			Env:  viper.GetString("APP_ENV"),
			Name: viper.GetString("APP_NAME"),
			Port: viper.GetString("APP_PORT"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
			Timezone: viper.GetString("DB_TIMEZONE"),
		},
		JWT: JWTConfig{
			Secret:      viper.GetString("JWT_SECRET"),
			ExpiredTime: jwtExpired,
		},
		Pagination: PaginationConfig{
			DefaultPageSize: viper.GetInt("DEFAULT_PAGE_SIZE"),
			MaxPageSize:     viper.GetInt("MAX_PAGE_SIZE"),
		},
	}

	if config.Database.User == "" || config.Database.Password == "" || config.Database.Name == "" {
		return nil, fmt.Errorf("database credentials are not fully set")
	}

	if config.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT secret is not set")
	}

	return config, nil
}
