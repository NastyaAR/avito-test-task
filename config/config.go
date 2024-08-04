package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Logger `yaml:"logger"`
	Db     `yaml:"postgres"`
}

type Logger struct {
	LogLevel string `yaml:"log-level"`
	LogFile  string `yaml:"log-file"`
}

type Db struct {
	Host         string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user" env:"POSTGRES_USER" env-default:"postgres"`
	Password     string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"password"`
	Db           string `yaml:"postgres-db" env:"POSTGRES_DB" env-default:"postgres"`
	DbTimeoutSec int    `yaml:"db-timeout-sec"`
}

func ReadConfig() (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig("config/config.yml", &cfg)
	if err != nil {
		return nil, fmt.Errorf("read config error: %v", err.Error)
	}
	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("read env error: %v", err.Error)
	}

	fmt.Println(cfg)
	return &cfg, nil
}
