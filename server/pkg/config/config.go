package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	ConnectionString string `yaml:"connectionString"`
}

type Config struct {
	DatabaseConfig DatabaseConfig `yaml:"database"`
}

func NewConfig(
	configPath string,
) (*Config, error) {
	var config Config
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}