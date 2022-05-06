package auth

import (
	"encoding/base64"
	"errors"
	"strings"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

type envoyIngressAuthAdapter struct {
	request *envoy_service_auth_v3.AttributeContext_HttpRequest
}

func NewEnvoyIngressAuthAdapter(req *envoy_service_auth_v3.AttributeContext_HttpRequest) AuthenticationCredentialProvider {
	return &envoyIngressAuthAdapter{request: req}
}

func (a *envoyIngressAuthAdapter) Credential() (AuthenticationCredential, error) {
	authHeader := a.request.Headers["authorization"]
	if authHeader == "" {
		return nil, errors.New("no authorization header found in request")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "basic") {
		return nil, errors.New("not a basic auth type")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	pair := strings.SplitN(string(decoded), ":", 2)
	if len(pair) != 2 || pair[0] == "" {
		return nil, errors.New("invalid basic auth decoded pair")
	}

	return &authCredential{userId: pair[0], userSecret: pair[1]}, nil
}
