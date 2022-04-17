package pdp

import (
	"context"
	"errors"
	"log"
	"os"

	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
	envoy_api_v3_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
)

type authorizationService struct {
	config       *common_config.Config
	policyEngine *PolicyEngine
}

func NewAuthorizationService(config *common_config.Config) (envoy_service_auth_v3.AuthorizationServer, error) {
	engine, err := NewPolicyEngine(os.Getenv("PDP_POLICY_PATH"), true)
	if err != nil {
		return &authorizationService{}, err
	}

	return &authorizationService{config: config, policyEngine: engine}, nil
}

func (s *authorizationService) Check(ctx context.Context,
	req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {

	httpReq := req.Attributes.Request.Http

	upstreamArtefact, err := s.resolveRequestedArtefact(httpReq)
	if err != nil {
		return &envoy_service_auth_v3.CheckResponse{}, err
	}

	log.Printf("Authorizing upstream req: [%s/%s/%s/%s][%s] %s", upstreamArtefact.Source.Type,
		upstreamArtefact.Group,
		upstreamArtefact.Name, upstreamArtefact.Version,
		httpReq.Method, httpReq.Path)

	return &envoy_service_auth_v3.CheckResponse{
		HttpResponse: &envoy_service_auth_v3.CheckResponse_OkResponse{
			OkResponse: &envoy_service_auth_v3.OkHttpResponse{
				Headers: []*envoy_api_v3_core.HeaderValueOption{
					{
						Append: &wrappers.BoolValue{Value: true},
						Header: &envoy_api_v3_core.HeaderValue{
							Key:   "x-pdp-authorized",
							Value: "true",
						},
					},
				},
			},
		},
		Status: &status.Status{
			Code: int32(code.Code_OK),
		},
	}, nil
}

func (s *authorizationService) resolveRequestedArtefact(req *envoy_service_auth_v3.AttributeContext_HttpRequest) (common_models.Artefact, error) {
	for _, upstream := range s.config.Global.Upstreams {
		if upstream.MatchPath(req.Path) {
			return upstream.Path2Artefact(req.Path)
		}
	}

	return common_models.Artefact{}, errors.New("failed to resolve artefact from upstream config")
}
