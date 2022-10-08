package config

import (
	"fmt"
	"os"

	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
)

type configHolder struct {
	configRepository ConfigRepository
}

var (
	globalConfig *configHolder = nil
)

func Bootstrap(file string, reloadOnChange bool) {
	if globalConfig != nil {
		panic("config is already bootstrapped")
	}

	if file == "" {
		file = os.Getenv("GLOBAL_CONFIG_PATH")
		if file == "" {
			panic("no config source available")
		}
	}

	repository, err := NewConfigFileRepository(file, false, reloadOnChange)
	if err != nil {
		panic(err)
	}

	globalConfig = &configHolder{
		configRepository: repository,
	}
}

// Returns the configuration data as per spec
func Get() *config_api.GatewayConfiguration {
	cfg, err := Current().configRepository.LoadGatewayConfiguration()
	if err != nil {
		panic(err)
	}

	return cfg
}

// Returns a wrapped configuration data with the wrapper
// providing some convenient utility method
func Current() *configHolder {
	if globalConfig == nil {
		panic("config is used without bootstrap")
	}

	return globalConfig
}

func (cfg *configHolder) GetMessagingConfigByName(name string) (*config_api.MessagingAdapter, error) {
	if mc, ok := Get().Messaging[name]; ok {
		return mc, nil
	} else {
		return nil, fmt.Errorf("messaging adapter not found with name: %s", name)
	}
}

func (cfg *configHolder) GetAuthenticatorByName(name string) (*config_api.GatewayAuthenticator, error) {
	if a, ok := Get().Authenticators[name]; ok {
		return a, nil
	} else {
		return nil, fmt.Errorf("authenticator not found with name: %s", name)
	}
}

func (cfg *configHolder) PdpServiceConfig() *config_api.PdpServiceConfig {
	return Get().Services.GetPdp()
}

func (cfg *configHolder) DcsServiceConfig() *config_api.DcsServiceConfig {
	return Get().Services.GetDcs()
}

func (cfg *configHolder) Upstreams() []*config_api.GatewayUpstream {
	return Get().Upstreams
}
