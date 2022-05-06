package auth

import (
	"errors"
	"fmt"

	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
)

// Implement basic auth for gateway ingress
type basicAuthProvider struct {
	config      common_config.AuthenticatorConfig
	file        string
	credentials map[string]string
}

func NewIngressBasicAuthService(config common_config.AuthenticatorConfig) (IngressAuthenticationService, error) {
	// TODO: read and parse the htpasswd file
	return &basicAuthProvider{config: config}, nil
}

func (p *basicAuthProvider) Authenticate(cp AuthenticationCredentialProvider) (AuthenticatedIdentity, error) {
	creds, err := cp.Credential()
	if err != nil {
		return nil, errors.New("no credential found")
	}

	fmt.Printf("creds: %p\n", &creds)
	return nil, errors.New("unimplemented")
}
