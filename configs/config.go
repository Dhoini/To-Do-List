package configs

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Db struct {
		Dsn string `mapstructure:"DSN"`
	} `mapstructure:"DB"`
	Auth struct {
		Secret        string        `mapstructure:"SECRET"`
		TokenLifetime time.Duration `mapstructure:"TOKEN_LIFETIME"`
	} `mapstructure:"AUTH"`
	Server struct {
		Port         int           `mapstructure:"PORT"`
		ReadTimeout  time.Duration `mapstructure:"READ_TIMEOUT"`
		WriteTimeout time.Duration `mapstructure:"WRITE_TIMEOUT"`
	} `mapstructure:"SERVER"`
	RateLimit struct {
		MaxRequests float64       `mapstructure:"MAX_REQUESTS"`
		Burst       int           `mapstructure:"BURST"`
		TTL         time.Duration `mapstructure:"TTL"`
	} `mapstructure:"RATE_LIMIT"`
}

// LoadConfig загружает конфигурацию из файла config.yaml/config.json и переменных окружения
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // Имя файла конфигурации (без расширения)
	viper.SetConfigType("yaml")   //  yaml,  json, toml
	viper.AddConfigPath(".")      // Искать в текущей директории
	viper.AutomaticEnv()          // Автоматически читать переменные окружения

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil { // viper использует mapstructure
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// Валидация
	if config.Db.Dsn == "" {
		return nil, fmt.Errorf("database DSN is required")
	}
	if config.Auth.Secret == "" {
		return nil, fmt.Errorf("auth secret is required")
	}
	if config.Server.Port == 0 {
		return nil, fmt.Errorf("server port is required")
	}
	if config.Auth.TokenLifetime == 0 { // Добавляем валидацию TokenLifetime
		config.Auth.TokenLifetime = time.Hour * 24 // Значение по умолчанию
	}

	return &config, nil
}
