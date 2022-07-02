package config

import (
	"errors"
	"os"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
	"gopkg.in/yaml.v2"
)

type MessagingConfig struct {
	Url string `json:"url"`
}

type EventPublisherConfig struct {
	TopicMappings map[string]string `yaml:"topics"`
}

type TapServiceConfig struct {
	Publisher EventPublisherConfig `yaml:"publisher"`
}

type PdsClientConfig struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	UseMtls bool   `yaml:"mtls"`
	Type    string `yaml:"type"`
}

type PdpServiceConfig struct {
	MonitorMode bool                 `yaml:"monitor_mode"`
	Publisher   EventPublisherConfig `yaml:"publisher"`
	PdsClient   PdsClientConfig      `yaml:"pds_client"`
}

type DcsServiceConfig struct {
	Publisher EventPublisherConfig `yaml:"publisher"`
}

type SecretConfig struct {
	Source string `yaml:"source"`
	Key    string `yaml:"key"`
}

type AuthenticatorConfig struct {
	Type   string            `yaml:"type"`
	Params map[string]string `yaml:"params"`
}

type GlobalConfig struct {
	Upstreams      []models.ArtefactUpStream      `yaml:"upstreams"`
	Messaging      MessagingConfig                `yaml:"messaging"`
	TapService     TapServiceConfig               `yaml:"tap"`
	PdpService     PdpServiceConfig               `yaml:"pdp"`
	DcsService     DcsServiceConfig               `yaml:"dcs"`
	Secrets        map[string]SecretConfig        `yaml:"secrets"`
	Authenticators map[string]AuthenticatorConfig `yaml:"authenticators"`
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
