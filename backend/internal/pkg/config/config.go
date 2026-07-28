package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func newReplacer() *strings.Replacer {
	return strings.NewReplacer(".", "_")
}

type Config struct {
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Server     ServerConfig     `mapstructure:"server"`
	SMTP       SMTPConfig       `mapstructure:"smtp"`
	Storage    StorageConfig    `mapstructure:"storage"`
	CORS       CORSConfig       `mapstructure:"cors"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	Metrics    MetricsConfig    `mapstructure:"metrics"`
	OTEL       OTELConfig       `mapstructure:"otel"`
}

type DatabaseConfig struct {
	URL         string `mapstructure:"url"`
	MaxConns    int32  `mapstructure:"max_conns"`
	MinConns    int32  `mapstructure:"min_conns"`
}

type RedisConfig struct {
	URL      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret        string        `mapstructure:"secret"`
	RefreshSecret string        `mapstructure:"refresh_secret"`
	AccessTTL     time.Duration `mapstructure:"access_ttl"`
	RefreshTTL    time.Duration `mapstructure:"refresh_ttl"`
}

type ServerConfig struct {
	Port  int    `mapstructure:"port"`
	Mode  string `mapstructure:"mode"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type StorageConfig struct {
	Path string `mapstructure:"path"`
	Type string `mapstructure:"type"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type MetricsConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

type OTELConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Endpoint string `mapstructure:"endpoint"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetConfigFile("config.yaml")
	v.AutomaticEnv()
	v.SetEnvPrefix("DEVFLOW")
	v.SetEnvKeyReplacer(newReplacer())

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	v.SetDefault("database.max_conns", 25)
	v.SetDefault("database.min_conns", 5)
	v.SetDefault("jwt.access_ttl", "15m")
	v.SetDefault("jwt.refresh_ttl", "168h")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "release")
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.format", "json")
	v.SetDefault("metrics.enabled", true)
	v.SetDefault("storage.type", "local")
	v.SetDefault("storage.path", "./data/storage")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}
	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.Redis.URL == "" {
		return fmt.Errorf("REDIS_URL is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.JWT.RefreshSecret == "" {
		return fmt.Errorf("JWT_REFRESH_SECRET is required")
	}
	return nil
}
