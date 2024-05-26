package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	version string = "1.0"
)

type Config struct {
	Project  Project  `yaml:"project"`
	Telegram Telegram `yaml:"telegram"`
	DBConfig DBConfig `yaml:"database"`
	Server   Server   `yaml:"server"`
}

type Project struct {
	Name        string `yaml:"name"`
	ServiceName string `yaml:"serviceName"`
	Version     string
}

type Server struct {
	Port  int    `yaml:"port"`
	Token string `yaml:"token"`
}

type Telegram struct {
	Token string `yaml:"token" envconfig:"TELEGRAM_TOKEN" validate:"required"`
}

type DBConfig struct {
	Host         string `yaml:"host" envconfig:"DB_HOST"`
	Port         int    `yaml:"port" envconfig:"DB_PORT"`
	DatabaseName string `yaml:"database_name" envconfig:"DB_DATABASE_NAME"`
	Username     string `yaml:"username" envconfig:"DB_USERNAME"`
	Password     string `yaml:"password" envconfig:"DB_PASSWORD"`
}

func New() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "./configs/config.yaml"
	}

	config := &Config{}

	if err := fromYaml(path, config); err != nil {
		fmt.Printf("couldn'n load config from %s: %s\r\n", path, err.Error())
	}

	if err := fromEnv(config); err != nil {
		fmt.Printf("couldn'n load config from env: %s\r\n", err.Error())
	}

	if err := validate(config); err != nil {
		return nil, err
	}

	config.Project.Version = version

	return config, nil
}

func fromYaml(path string, config *Config) error {
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

func fromEnv(config *Config) error {
	return envconfig.Process("", config)
}

func validate(cfg *Config) error {
	if cfg.Telegram.Token == "" {
		return fmt.Errorf("config: %s is not set", "TELEGRAM_TOKEN")
	}

	if cfg.DBConfig.DatabaseName == "" {
		return fmt.Errorf("config: %s is not set", "DB_DATABASE_NAME")
	}

	if cfg.DBConfig.Username == "" {
		return fmt.Errorf("config: %s is not set", "DB_USERNAME")
	}

	return nil
}
