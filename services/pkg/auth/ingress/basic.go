package ingress

import (
	"errors"
	"fmt"

	"github.com/abhisek/supply-chain-gateway/services/pkg/auth"
)

// Implement basic auth for gateway ingress
type basicAuthProvider struct {
	file        string
	credentials map[string]string
}

func NewIngressBasicAuthService() (auth.IngressAuthenticationService, error) {
	// TODO: read and parse the htpasswd file
	return &basicAuthProvider{}, nil
}

func (p *basicAuthProvider) Authenticate(cp auth.AuthenticationCredentialProvider) (auth.AuthenticatedIdentity, error) {
	creds, err := cp.Credential()
	if err != nil {
		return nil, errors.New("no credential found")
	}

	fmt.Printf("creds: %p\n", &creds)
	return nil, errors.New("unimplemented")
}
