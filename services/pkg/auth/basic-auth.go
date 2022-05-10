package auth

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"golang.org/x/crypto/bcrypt"
)

var (
	basicAuthUserNotFound        = errors.New("user not found in basic auth db")
	basicAuthFailed              = errors.New("authentication denied")
	basicAuthCredentialNotFound  = errors.New("credential not found in request")
	basicAuthHashTypeUnsupported = errors.New("hash type is not supported")
)

// Implement basic auth for gateway ingress
type basicAuthProvider struct {
	config      common_config.AuthenticatorConfig
	file        string
	credentials map[string]string
}

func NewIngressBasicAuthService(config common_config.AuthenticatorConfig) (IngressAuthenticationService, error) {
	p := &basicAuthProvider{config: config, file: config.Params["htpasswd_file"]}
	if err := p.loadCredentials(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *basicAuthProvider) Authenticate(ctx context.Context, cp AuthenticationCredentialProvider) (AuthenticatedIdentity, error) {
	creds, err := cp.Credential()
	if err != nil {
		return nil, basicAuthCredentialNotFound
	}

	hp, ok := p.credentials[creds.UserId()]
	if !ok {
		return nil, basicAuthUserNotFound
	}

	err = p.safeCompareHash(creds.UserSecret(), hp)
	if err != nil {
		return nil, err
	}

	return NewAuthIdentity(AuthIdentityTypeBasicAuth,
		creds.UserId(),
		fmt.Sprintf("Basic Auth User: %s", creds.UserId())), nil

}

func (p *basicAuthProvider) loadCredentials() error {
	log.Printf("Loading basic auth credentials from: %s", p.file)

	file, err := os.OpenFile(p.file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	s := make(map[string]string, 0)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ":", 2)
		if len(parts) == 2 {
			s[parts[0]] = parts[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	p.credentials = s
	return nil
}

func (p *basicAuthProvider) safeCompareHash(password string, hash string) error {
	if !strings.HasPrefix(hash, "$2y$") {
		return basicAuthHashTypeUnsupported
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return basicAuthFailed
	}

	return nil
}
