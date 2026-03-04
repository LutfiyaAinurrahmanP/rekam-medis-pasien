package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig
	Database   DatabaseConfig
	JWT        JWTConfig
	Pagination PaginationConfig
	Redis      RedisConfig
	Kafka      KafkaConfig
}

// KafkaConfig menyimpan konfigurasi koneksi ke Apache Kafka.
type KafkaConfig struct {
	// Brokers adalah daftar alamat Kafka broker yang dipisahkan koma,
	// contoh: "localhost:9092,kafka2:9092"
	Brokers  []string
	ClientID string
	// Enabled menentukan apakah Kafka diaktifkan. Jika false,
	// event publishing dinonaktifkan dan aplikasi tetap berjalan.
	Enabled  bool
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

type RedisConfig struct {
	Host       string
	Port       string
	Password   string
	DB         int
	DefaultTTL time.Duration
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
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("REDIS_TTL", "5m")
	viper.SetDefault("KAFKA_BROKERS", "localhost:9092")
	viper.SetDefault("KAFKA_CLIENT_ID", "sirekam-medis")
	viper.SetDefault("KAFKA_ENABLED", true)

	jwtExpired, err := time.ParseDuration(viper.GetString("JWT_EXPIRED_TIME"))
	if err != nil {
		jwtExpired = 24 * time.Hour
	}

	redisTTL, err := time.ParseDuration(viper.GetString("REDIS_TTL"))
	if err != nil {
		redisTTL = 5 * time.Minute
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
		Redis: RedisConfig{
			Host:       viper.GetString("REDIS_HOST"),
			Port:       viper.GetString("REDIS_PORT"),
			Password:   viper.GetString("REDIS_PASSWORD"),
			DB:         viper.GetInt("REDIS_DB"),
			DefaultTTL: redisTTL,
		},
		Kafka: KafkaConfig{
			Brokers:  strings.Split(viper.GetString("KAFKA_BROKERS"), ","),
			ClientID: viper.GetString("KAFKA_CLIENT_ID"),
			Enabled:  viper.GetBool("KAFKA_ENABLED"),
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
