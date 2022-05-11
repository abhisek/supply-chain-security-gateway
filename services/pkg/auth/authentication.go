package auth

/**
	[Gateway] -> [Ingress Auth] -> PDP
	[Gateway] -> TAP -> [Egress Auth] -> Upstream Repo

	Risk?

	User tricking gateway to send credentials to malicious user controlled endpoint
**/

import (
	"context"

	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
)

const (
	// PDP will lookup ingress authenticators
	AuthStageIngress = "ingress" // Gateway Auth

	// Tap will lookup egress authenticators
	AuthStageEgress = "egress" // Upstream Auth

	AuthTypeNoAuth = "noauth"
	AuthTypeBasic  = "basic"
	AuthTypeOIDC   = "oidc"

	AuthIdentityTypeBasicAuth = "BasicAuth"
)

type AuthenticationProvider interface {
	IngressAuthService(common_models.ArtefactUpStream) (IngressAuthenticationService, error)
	EgressAuthService(common_models.ArtefactRepository) (EgressAuthenticationService, error)
}

// Adapter to wrap Envoy request to get credentials
type AuthenticationCredentialProvider interface {
	Credential() (AuthenticationCredential, error)
}

// A provided or obtained credential for authentication
type AuthenticationCredential interface {
	UserId() string
	UserSecret() string
}

// Authenticated identity used in Ingress auth
type AuthenticatedIdentity interface {
	Type() string
	Id() string
	Name() string
}

// Authentication for gateway users
type IngressAuthenticationService interface {
	Authenticate(context.Context, AuthenticationCredentialProvider) (AuthenticatedIdentity, error)
}

// Apply credentials to outgoing request to repo
type AuthenticationCredentialApplier interface {
	Apply(AuthenticationCredential) error
}

// Authenticate upstream repo request
type EgressAuthenticationService interface {
	Authenticate(AuthenticationCredentialApplier) error
}
