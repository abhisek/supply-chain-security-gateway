package auth

import (
	"errors"
	"fmt"

	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
)

type authProvider struct {
	// Unbounded cache, should not be a problem because the
	// number of providers can be limited
	ingressCache map[string]IngressAuthenticationService
	egressCache  map[string]EgressAuthenticationService
}

func NewAuthenticationProvider() AuthenticationProvider {
	return &authProvider{}
}

func (a *authProvider) IngressAuthService(upstream common_models.ArtefactUpStream) (IngressAuthenticationService, error) {
	cf := func(s func(c *config_api.GatewayAuthenticator) (IngressAuthenticationService, error)) (IngressAuthenticationService, error) {
		cfg, err := config.GetAuthenticatorByName(upstream.Authentication.Provider)
		if err != nil {
			return nil, err
		}

		return s(cfg)
	}

	// TODO: Implement a cache for services to prevent reinitialize the same
	// authenticator, uniquely identified by a name

	switch upstream.Authentication.Type {
	case AuthTypeNoAuth:
		return NewIngressNoAuthService()
	case AuthTypeBasic:
		return cf(func(c *config_api.GatewayAuthenticator) (IngressAuthenticationService, error) {
			return NewIngressBasicAuthService(c.GetBasicAuth())
		})
	default:
		return nil, fmt.Errorf("no auth service available for: %s", upstream.Authentication.Provider)
	}
}

func (a *authProvider) EgressAuthService(common_models.ArtefactRepository) (EgressAuthenticationService, error) {
	return nil, errors.New("unimplemented")
}
