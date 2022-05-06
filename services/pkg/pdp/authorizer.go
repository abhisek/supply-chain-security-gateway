package pdp

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	pds_api "github.com/abhisek/supply-chain-gateway/services/gen"
	"github.com/abhisek/supply-chain-gateway/services/pkg/auth"
	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
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
	config            *common_config.Config
	authProvider      auth.AuthenticationProvider
	policyEngine      *PolicyEngine
	policyDataService pds_api.PolicyDataServiceClient
}

func NewAuthorizationService(config *common_config.Config, p pds_api.PolicyDataServiceClient) (envoy_service_auth_v3.AuthorizationServer, error) {
	engine, err := NewPolicyEngine(os.Getenv("PDP_POLICY_PATH"), true)
	if err != nil {
		return &authorizationService{}, err
	}

	authProvider := auth.NewAuthenticationProvider(config)
	return &authorizationService{config: config,
		authProvider: authProvider,
		policyEngine: engine, policyDataService: p}, nil
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

	nctx, ncancel := context.WithTimeout(ctx, 2*time.Second)
	defer ncancel()
	resp, err := s.policyDataService.FindVulnerabilitiesByArtefact(nctx, &pds_api.FindVulnerabilityByArtefactRequest{
		Artefact: &pds_api.Artefact{
			Ecosystem: upstreamArtefact.OpenSsfEcosystem(),
			Group:     upstreamArtefact.Group,
			Name:      upstreamArtefact.Name,
			Version:   upstreamArtefact.Version,
		},
	})

	var vulnerabilities []common_models.ArtefactVulnerability
	if err != nil {
		log.Printf("Failed to enrich artefact with vulnerability information: %v", err)
	} else {
		log.Printf("Enriched artefact (%s/%s/%s) with vulnerabilities: %s",
			upstreamArtefact.Group, upstreamArtefact.Name, upstreamArtefact.Version,
			utils.Introspect(resp.Vulnerabilities))

		vulnerabilities = s.mapVulnerabilities(resp.Vulnerabilities)
	}

	log.Printf("Authorizing upstream req from %s: [%s/%s/%s/%s][%s] %s",
		userId,
		upstreamArtefact.Source.Type,
		upstreamArtefact.Group,
		upstreamArtefact.Name, upstreamArtefact.Version,
		httpReq.Method, httpReq.Path)

	policyRespose, err := s.policyEngine.Evaluate(NewPolicyInput(upstreamArtefact, upstream, vulnerabilities))
	if err != nil {
		log.Printf("Failed to evaluate policy: %s", err.Error())
		return &envoy_service_auth_v3.CheckResponse{}, err
	}

	if !s.config.Global.PdpService.MonitorMode && !policyRespose.Allowed() {
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

	authService, err := s.authProvider.IngressAuthService(upstream)
	if err != nil {
		return "", err
	}

	identity, err := authService.Authenticate(auth.NewEnvoyIngressAuthAdapter(req))
	if err != nil {
		return "", err
	}

	return identity.Id(), nil
}

// Convert pds_api.VulnerabilityMeta to common_models.ArtefactVulnerability
func (s *authorizationService) mapVulnerabilities(src []*pds_api.VulnerabilityMeta) []common_models.ArtefactVulnerability {
	target := []common_models.ArtefactVulnerability{}

	if src == nil || len(src) == 0 {
		return target
	}

	for _, s := range src {
		mv := common_models.ArtefactVulnerability{
			Name: s.Title,
			Id: common_models.ArtefactVulnerabilityId{
				Source: s.Source,
				Id:     s.Id,
			},
			Scores: []common_models.ArtefactVulnerabilityScore{},
		}

		switch s.Severity {
		case pds_api.VulnerabilitySeverity_CRITICAL:
			mv.Severity = common_models.ArtefactVulnerabilitySeverityCritical
			break
		case pds_api.VulnerabilitySeverity_HIGH:
			mv.Severity = common_models.ArtefactVulnerabilitySeverityHigh
			break
		case pds_api.VulnerabilitySeverity_MEDIUM:
			mv.Severity = common_models.ArtefactVulnerabilitySeverityMedium
			break
		case pds_api.VulnerabilitySeverity_LOW:
			mv.Severity = common_models.ArtefactVulnerabilitySeverityLow
		default:
			mv.Severity = common_models.ArtefactVulnerabilitySeverityInfo
		}

		for _, score := range s.Scores {
			mv.Scores = append(mv.Scores, common_models.ArtefactVulnerabilityScore{
				Type:  score.Type,
				Value: score.Value,
			})
		}

		target = append(target, mv)
	}

	return target
}
