package auth

type noAuthProvider struct{}

func NewIngressNoAuthService() (IngressAuthenticationService, error) {
	return &noAuthProvider{}, nil
}

func (p *noAuthProvider) Authenticate(cp AuthenticationCredentialProvider) (AuthenticatedIdentity, error) {
	return NewAuthIdentity("", "No Auth", "No Auth"), nil
}
