package config

import (
	"errors"
	"os"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
	"gopkg.in/yaml.v2"
)

type GlobalConfig struct {
	Upstreams []models.ArtefactUpStream `yaml:"upstreams"`
}

type Config struct {
	Global GlobalConfig
}

func LoadGlobal(file string) (*Config, error) {
	if file == "" {
		file = os.Getenv("GLOBAL_CONFIG_PATH")
		if file == "" {
			return &Config{}, errors.New("failed to find config path")
		}
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return &Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config.Global)
	if err != nil {
		return &Config{}, err
	}

	return &config, nil
}
