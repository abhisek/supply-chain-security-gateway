package config

import (
	"fmt"
	"os"
	"sync"

	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
	"gopkg.in/yaml.v2"
)

type configFileRepository struct {
	path                 string
	gatewayConfiguration *config_api.GatewayConfiguration
	m                    sync.RWMutex
}

func NewConfigFileRepository(path string) (ConfigRepository, error) {
	r := &configFileRepository{path: path}

	err := r.load()
	if err == nil {
		err = r.monitorForChange()
	}

	return r, err
}

func (c *configFileRepository) LoadGatewayConfiguration() (*config_api.GatewayConfiguration, error) {
	var err error = nil
	if c.gatewayConfiguration != nil {
		err = fmt.Errorf("gateway configuration is not loaded")
	}

	c.m.RLock()
	defer c.m.RUnlock()

	return c.gatewayConfiguration, err
}

func (c *configFileRepository) SaveGatewayConfiguration(config *config_api.GatewayConfiguration) error {
	return fmt.Errorf("persisting gateway configuration is not supported")
}

func (c *configFileRepository) load() error {
	file, err := os.Open(c.path)
	if err != nil {
		return err
	}

	var gatewayConfiguration config_api.GatewayConfiguration
	err = yaml.NewDecoder(file).Decode(&gatewayConfiguration)
	if err != nil {
		return err
	}

	err = gatewayConfiguration.Validate()
	if err != nil {
		return err
	}

	c.m.Lock()
	defer c.m.Unlock()

	c.gatewayConfiguration = &gatewayConfiguration
	return nil
}

func (c *configFileRepository) monitorForChange() error {
	return nil
}
