package pdp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abhisek/supply-chain-gateway/services/pkg/auth"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
	envoy_api_v3_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
)

func (s *authorizationService) authenticateForUpstream(ctx context.Context,
	upstream common_models.ArtefactUpStream,
	req *envoy_service_auth_v3.AttributeContext_HttpRequest) (auth.AuthenticatedIdentity, error) {
	if !upstream.NeedAuthentication() {
		return auth.AnonymousIdentity(), nil
	}

	if req.Method == "HEAD" {
		return auth.AnonymousIdentity(), nil
	}

	authService, err := s.authProvider.IngressAuthService(upstream)
	if err != nil {
		return nil, err
	}

	identity, err := authService.Authenticate(ctx, auth.NewEnvoyIngressAuthAdapter(req))
	if err != nil {
		return nil, err
	}

	return identity, nil
}

func (s *authorizationService) authenticationChallenge(ctx context.Context,
	upstream common_models.ArtefactUpStream,
	req *envoy_service_auth_v3.AttributeContext_HttpRequest) (*envoy_service_auth_v3.CheckResponse, error) {

	authChallenge := fmt.Sprintf("Basic realm=\"Authentication for upstream %s at %s\"",
		upstream.Name, upstream.RoutingRule.Prefix)

	return &envoy_service_auth_v3.CheckResponse{
		HttpResponse: &envoy_service_auth_v3.CheckResponse_DeniedResponse{
			DeniedResponse: &envoy_service_auth_v3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{
					Code: http.StatusUnauthorized,
				},
				Headers: []*envoy_api_v3_core.HeaderValueOption{
					{
						Append: &wrappers.BoolValue{Value: false},
						Header: &envoy_api_v3_core.HeaderValue{
							Key:   "WWW-Authenticate",
							Value: authChallenge,
						},
					},
				},
			},
		},
		Status: &status.Status{
			Code: int32(code.Code_UNAUTHENTICATED),
		},
	}, nil
}
