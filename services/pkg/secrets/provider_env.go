package secrets

import (
	"fmt"
	"os"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
)

type envSecretProvider struct{}

func NewEnvSecretProvider() (SecretProvider, error) {
	return &envSecretProvider{}, nil
}

func (p *envSecretProvider) GetSecret(name string) (string, error) {
	v, b := os.LookupEnv(name)
	if !b || utils.IsEmptyString(v) {
		return "", fmt.Errorf("secret %s returned empty or non-existent", name)
	}

	return v, nil
}
