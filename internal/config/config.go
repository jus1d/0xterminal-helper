package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	EnvLocal       = "local"
	EnvDevelopment = "dev"
	EnvProduction  = "prod"
)

// Config combine all sub-configs structures
type Config struct {
	Env      string   `yaml:"env" env-required:"true"`
	Telegram Telegram `yaml:"telegram"`
	Postgres Postgres `yaml:"postgres"`
	OCR      OCR      `yaml:"ocr"`
}

// Telegram represents structure with credentials for Telegram bot connection
type Telegram struct {
	Token string `yaml:"token"`
}

// Pstgres represents structure with credentials for PostgreSQL database
type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	ModeSSL  string `yaml:"sslmode"`
}

// OCR represents structure with credentials for OCR service (ocr.space currently)
type OCR struct {
	Token string `yaml:"token"`
}

// MustLoad loads config to a new Config instance and return it's pointer.
func MustLoad() *Config {
	_ = godotenv.Load()

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatalf("missed CONFIG_PATH parameter")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist at: %s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &config
}
