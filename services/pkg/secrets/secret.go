package secrets

import (
	"errors"
	"fmt"
	"log"

	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
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

func getSecret(secret common_config.SecretConfig, evictCache bool) (string, error) {
	if evictCache {
		// TODO: Evict the cache
	}

	// TODO: Lookup cache

	a, err := resolveProvider(secret.Source)
	if err != nil {
		return "", err
	}

	return a.GetSecret(secret.Key)
}

func resolveProvider(name string) (SecretProvider, error) {
	switch name {
	case SecretProviderEnv:
		return NewEnvSecretProvider()
	default:
		return nil, fmt.Errorf("provider not found for %s", name)
	}
}

func GetSecret(secret common_config.SecretConfig) (string, error) {
	return getSecret(secret, false)
}

func Resolve(id string, config *common_config.Config) (common_config.SecretConfig, error) {
	return common_config.SecretConfig{}, errors.New("unimplemented")
}
