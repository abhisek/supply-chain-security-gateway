package config

/**
The config module provides access to currently available configuration.
A repository is used as the source of configuration.

Configuration at a high level is:

1. Bootstrap (can be dynamically updated)
2. Contextual (depends on current context data)
*/

import (
	"fmt"
	"os"

	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
)

type configHolder struct {
	ctx              any // TODO: Standardize domain context
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

// Get a contextual config
func Contextual(ctx any) *configHolder {
	return current().withContext(ctx)
}

// Returns a wrapped configuration data with the wrapper
// providing some convenient utility method
func current() *configHolder {
	if globalConfig == nil {
		panic("config is used without bootstrap")
	}

	return globalConfig
}

func (cfg *configHolder) withContext(ctx any) *configHolder {
	return &configHolder{
		configRepository: cfg.configRepository,
		ctx:              ctx,
	}
}

// Returns the configuration data as per spec
func g() *config_api.GatewayConfiguration {
	cfg, err := current().configRepository.LoadGatewayConfiguration()
	if err != nil {
		panic(err)
	}

	return cfg
}

func GetMessagingConfigByName(name string) (*config_api.MessagingAdapter, error) {
	if mc, ok := g().Messaging[name]; ok {
		return mc, nil
	} else {
		return nil, fmt.Errorf("messaging adapter not found with name: %s", name)
	}
}

func GetAuthenticatorByName(name string) (*config_api.GatewayAuthenticator, error) {
	if a, ok := g().Authenticators[name]; ok {
		return a, nil
	} else {
		return nil, fmt.Errorf("authenticator not found with name: %s", name)
	}
}

func PdpServiceConfig() *config_api.PdpServiceConfig {
	return g().Services.GetPdp()
}

func DcsServiceConfig() *config_api.DcsServiceConfig {
	return g().Services.GetDcs()
}

func TapServiceConfig() *config_api.TapServiceConfig {
	return g().Services.GetTap()
}

func Upstreams() []*config_api.GatewayUpstream {
	return g().GetUpstreams()
}

func GetSecret(name string) (*config_api.GatewaySecret, error) {
	if s, ok := g().Secrets[name]; ok {
		return s, nil
	} else {
		return nil, fmt.Errorf("no secret found with name: %s", name)
	}
}
