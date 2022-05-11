package auth

import (
	"errors"
	"fmt"

	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
)

type authProvider struct {
	config *common_config.Config
}

func NewAuthenticationProvider(config *common_config.Config) AuthenticationProvider {
	return &authProvider{config: config}
}

func (a *authProvider) IngressAuthService(upstream common_models.ArtefactUpStream) (IngressAuthenticationService, error) {
	cf := func(s func(c common_config.AuthenticatorConfig) (IngressAuthenticationService, error)) (IngressAuthenticationService, error) {
		cfg, ok := a.config.Global.Authenticators[upstream.Authentication.Provider]
		if !ok {
			return nil, fmt.Errorf("no authenticator defined for: %s", upstream.Authentication.Provider)
		}

		return s(cfg)
	}

	switch upstream.Authentication.Type {
	case AuthTypeNoAuth:
		return NewIngressNoAuthService()
	case AuthTypeBasic:
		return cf(func(c common_config.AuthenticatorConfig) (IngressAuthenticationService, error) {
			return NewIngressBasicAuthService(c)
		})
	default:
		return nil, fmt.Errorf("no auth service available for: %s", upstream.Authentication.Provider)
	}
}

func (a *authProvider) EgressAuthService(common_models.ArtefactRepository) (EgressAuthenticationService, error) {
	return nil, errors.New("unimplemented")
}
