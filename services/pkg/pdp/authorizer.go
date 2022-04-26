package pdp

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"os"
	"strings"

	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
	envoy_api_v3_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
)

var (
	errPolicyDeniedUpStreamRequest = errors.New("policy denied upstream request")
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

	upstreamArtefact, upstream, err := s.resolveRequestedArtefact(httpReq)
	if err != nil {
		log.Printf("No artefact resolved: %s", err.Error())
		return &envoy_service_auth_v3.CheckResponse{}, err
	}

	userId, err := s.authenticateForUpstream(upstream, httpReq)
	if err != nil {
		log.Printf("Error resolving userId: %v", err)
		return &envoy_service_auth_v3.CheckResponse{}, err
	}

	log.Printf("Authorizing upstream req from %s: [%s/%s/%s/%s][%s] %s",
		userId,
		upstreamArtefact.Source.Type,
		upstreamArtefact.Group,
		upstreamArtefact.Name, upstreamArtefact.Version,
		httpReq.Method, httpReq.Path)

	policyRespose, err := s.policyEngine.Evaluate(NewPolicyInputWithArtefact(upstreamArtefact, upstream))
	if err != nil {
		log.Printf("Failed to evaluate policy: %s", err.Error())
		return &envoy_service_auth_v3.CheckResponse{}, err
	}

	if !policyRespose.Allowed() {
		log.Printf("Policy denied upstream request")
		return &envoy_service_auth_v3.CheckResponse{}, errPolicyDeniedUpStreamRequest
	}

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

func (s *authorizationService) resolveRequestedArtefact(req *envoy_service_auth_v3.AttributeContext_HttpRequest) (common_models.Artefact,
	common_models.ArtefactUpStream, error) {
	for _, upstream := range s.config.Global.Upstreams {
		if upstream.MatchPath(req.Path) {
			a, err := upstream.Path2Artefact(req.Path)
			return a, upstream, err
		}
	}

	return common_models.Artefact{},
		common_models.ArtefactUpStream{},
		errors.New("failed to resolve artefact from upstream config")
}

// POC implementation of extracting UserId from basic auth header. Auth needs to be a
// service of its own with pluggable IDP support e.g. Github OIDC Token as password
// This helps us identify who is accessing the artefact so that violations can be attributed
func (s *authorizationService) authenticateForUpstream(upstream common_models.ArtefactUpStream,
	req *envoy_service_auth_v3.AttributeContext_HttpRequest) (string, error) {
	if !upstream.NeedAuthentication() {
		return "anonymous-upstream", nil
	}

	if req.Method == "HEAD" {
		return "anonymous-head", nil
	}

	authHeader := req.Headers["authorization"]
	if authHeader == "" {
		return "", errors.New("no authorization header found in request")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "basic") {
		return "", errors.New("not a basic auth type")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	pair := strings.SplitN(string(decoded), ":", 2)
	if len(pair) != 2 || pair[0] == "" {
		return "", errors.New("invalid basic auth decoded pair")
	}

	return pair[0], nil
}
