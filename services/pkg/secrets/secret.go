package secrets

import (
	"fmt"
	"log"

	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
)

const (
	SecretProviderEnv               = "environment"
	SecretProviderAwsSecretsManager = "aws-secrets-manager"
	SecretProviderHashicorpVault    = "hashicorp-vault"
)

/**
Possible source of secret

Environment
Vault
AWS Secrets Manager

Some of these secrets provider need authentication of their own. The adapter need special
environmental config to be able to fetch secret
*/

func init() {
	initCache()
}

func initCache() {
	log.Printf("Initializing secrets cache")
}

func getSecret(key string, evictCache bool) (string, error) {
	s, err := config.GetSecret(key)
	if err != nil {
		return "", err
	}

	if evictCache {
		// TODO: Evict the cache
	}

	// TODO: Lookup cache

	a, err := resolveProvider(s.Source)
	if err != nil {
		return "", err
	}

	// FIXME: Refactor the secret provider interface
	return a.GetSecret(s.GetEnvironment().Key)
}

func resolveProvider(src config_api.GatewaySecretSource) (SecretProvider, error) {
	switch src {
	case config_api.GatewaySecretSource_Environment:
		return NewEnvSecretProvider()
	default:
		return nil, fmt.Errorf("provider not found for %s", src.String())
	}
}

func GetSecret(key string) (string, error) {
	return getSecret(key, false)
}
