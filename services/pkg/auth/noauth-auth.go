package auth

import "context"

type noAuthProvider struct{}

func NewIngressNoAuthService() (IngressAuthenticationService, error) {
	return &noAuthProvider{}, nil
}

func (p *noAuthProvider) Authenticate(ctx context.Context, cp AuthenticationCredentialProvider) (AuthenticatedIdentity, error) {
	return NewAuthIdentity("", "No Auth", "No Auth"), nil
}
