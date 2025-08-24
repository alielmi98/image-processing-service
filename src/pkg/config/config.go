package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Password PasswordConfig
	Cors     CorsConfig
	JWT      JWTConfig
	RabbitMQ RabbitMQConfig
}

type ServerConfig struct {
	InternalPort string
	ExternalPort string
	RunMode      string
	Domain       string
}

type PostgresConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type PasswordConfig struct {
	IncludeChars     bool
	IncludeDigits    bool
	MinLength        int
	MaxLength        int
	IncludeUppercase bool
	IncludeLowercase bool
}

type CorsConfig struct {
	AllowOrigins string
}

type JWTConfig struct {
	AccessTokenExpireDuration  time.Duration
	RefreshTokenExpireDuration time.Duration
	Secret                     string
	RefreshSecret              string
}

type RabbitMQConfig struct {
	Host                 string
	Port                 string
	User                 string
	Password             string
	VHost                string
	SSLMode              string
	MaxIdleConns         int
	MaxOpenConns         int
	ConnMaxLifetime      time.Duration
	ProcessingRoutingKey string
	ResultRoutingKey     string
	PrefetchCount        int
	ReconnectDelay       time.Duration
	MaxReconnectAttempts int
}

func GetConfig() *Config {
	cfgPath := getConfigPath(os.Getenv("APP_ENV"))
	v, err := LoadConfig(cfgPath, "yml")
	if err != nil {
		log.Printf("failed to load config from %s: %v", cfgPath, err)

		fallbackPath := "../../pkg/config/config-development"
		log.Printf("trying fallback config: %s", fallbackPath)

		v, err = LoadConfig(fallbackPath, "yml")
		if err != nil {
			log.Fatalf("Error loading fallback config %s: %v", fallbackPath, err)
		}
	}

	cfg, err := ParseConfig(v)
	envPort := os.Getenv("PORT")
	if envPort != "" {
		cfg.Server.ExternalPort = envPort
		log.Printf("Set external port from environment -> %s", cfg.Server.ExternalPort)
	} else {
		cfg.Server.ExternalPort = cfg.Server.InternalPort
		log.Printf("Set external port from environment -> %s", cfg.Server.ExternalPort)
	}
	if err != nil {
		log.Fatalf("Error in parse config %v", err)
	}

	return cfg
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Printf("Unable to parse config: %v", err)
		return nil, err
	}
	return &cfg, nil
}
func LoadConfig(filename string, fileType string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType(fileType)
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		log.Printf("Unable to read config: %v", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}
	return v, nil
}

func getConfigPath(env string) string {
	if env == "docker" {
		return "/app/pkg/config/config-docker"
	} else if env == "production" {
		return "/app/pkg/config/config-production"
	} else {
		return "../pkg/config/config-development"
	}
}
