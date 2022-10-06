package config

import (
	"fmt"
	"os"

	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
)

const (
	configRepositoryTypeFile = "file"
)

// Define a repository interface to get the current configuration
// the repository implementation can internally refresh / cache
// configuration as required
type ConfigRepository interface {
	LoadGatewayConfiguration() (*config_api.GatewayConfiguration, error)
	SaveGatewayConfiguration(config *config_api.GatewayConfiguration) error
}

func NewConfigRepository() (ConfigRepository, error) {
	cType := os.Getenv("BOOTSTRAP_CONFIGURATION_REPOSITORY_TYPE")
	switch cType {
	case configRepositoryTypeFile:
		return NewConfigFileRepository(os.Getenv("BOOTSTRAP_CONFIGURATION_REPOSITORY_PATH"), false, true)
	}

	return nil, fmt.Errorf("unknown config repository type: %s", cType)
}
