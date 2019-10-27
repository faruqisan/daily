package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type (
	// Config sturct contain config
	Config struct {
		DatabaseConfig DatabaseConfig `yaml:"database"`
		RedisConfig    RedisConfig    `yaml:"redis"`
	}

	// DatabaseConfig config define configuration for database
	DatabaseConfig struct {
		DSN string `yaml:"dsn"`
	}

	// RedisConfig define config for redis
	RedisConfig struct {
		Host string `yaml:"host"`
	}
)

var cfg *Config

// Get function
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	file, err := ioutil.ReadFile("files/config.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
