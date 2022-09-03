package pdp

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/abhisek/supply-chain-gateway/services/pkg/auth"
	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/messaging"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
	envoy_api_v3_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	grpc_err_status "google.golang.org/grpc/status"

	event_api "github.com/abhisek/supply-chain-gateway/services/gen"
)

var (
	errPolicyDeniedUpStreamRequest = errors.New("policy denied upstream request")
)

type authorizationService struct {
	config            *common_config.Config
	authProvider      auth.AuthenticationProvider
	policyEngine      *PolicyEngine
	policyDataService PolicyDataClientInterface
	messagingService  messaging.MessagingService
}

func NewAuthorizationService(config *common_config.Config, p PolicyDataClientInterface,
	m messaging.MessagingService) (envoy_service_auth_v3.AuthorizationServer, error) {
	engine, err := NewPolicyEngine(os.Getenv("PDP_POLICY_PATH"), true)
	if err != nil {
		return &authorizationService{}, err
	}

	authProvider := auth.NewAuthenticationProvider(config)
	return &authorizationService{config: config,
		authProvider:      authProvider,
		policyEngine:      engine,
		policyDataService: p,
		messagingService:  m}, nil
}

func (s *authorizationService) Check(ctx context.Context,
	req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {

	httpReq := req.Attributes.Request.Http

	upstreamArtefact, upstream, err := s.resolveRequestedArtefact(httpReq)
	if err != nil {
		log.Printf("No artefact resolved: %s", err.Error())
		return &envoy_service_auth_v3.CheckResponse{}, err
	}

	identity, err := s.authenticateForUpstream(ctx, upstream, httpReq)
	if err != nil {
		log.Printf("Error resolving userId: %v", err)
		return &envoy_service_auth_v3.CheckResponse{}, err
	}

	nctx, ncancel := context.WithTimeout(ctx, 2*time.Second)
	defer ncancel()
	pdsResponse, enrichmentErr := s.policyDataService.GetPackageMetaByVersion(nctx,
		upstreamArtefact.OpenSsfEcosystem(), upstreamArtefact.Group,
		upstreamArtefact.Name, upstreamArtefact.Version)

	if enrichmentErr != nil {
		log.Printf("Failed to enrich artefact with vulnerability information: %v", enrichmentErr)
	} else {
		log.Printf("Enriched artefact (%s/%s/%s) with data: %s",
			upstreamArtefact.Group, upstreamArtefact.Name, upstreamArtefact.Version,
			utils.Introspect(pdsResponse))
	}

	log.Printf("Authorizing upstream req from %s: [%s/%s/%s/%s][%s] %s",
		identity.Id(),
		upstreamArtefact.Source.Type,
		upstreamArtefact.Group,
		upstreamArtefact.Name, upstreamArtefact.Version,
		httpReq.Method, httpReq.Path)

	policyRespose, err := s.policyEngine.Evaluate(ctx,
		NewPolicyInput(upstreamArtefact, upstream, pdsResponse.Vulnerabilities, pdsResponse.Licenses))
	if err != nil {
		log.Printf("Failed to evaluate policy: %s", err.Error())
		return &envoy_service_auth_v3.CheckResponse{}, err
	}

	gatewayDeny := !s.config.Global.PdpService.MonitorMode && !policyRespose.Allowed()
	s.publishDecisionEvent(ctx, identity.Id(), pdsResponse, !gatewayDeny,
		s.config.Global.PdpService.MonitorMode,
		upstream, upstreamArtefact, policyRespose, enrichmentErr)

	if gatewayDeny {
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

// TODO - Refactor from using N args to using a builder
func (s *authorizationService) publishDecisionEvent(ctx context.Context, userId string,
	pdsResponse PolicyDataServiceResponse,
	gw_allowed bool, monitor_mode bool, upstream common_models.ArtefactUpStream,
	artefact common_models.Artefact, result PolicyResponse, enrichmentErr error) {

	eh := common_models.NewSpecHeaderWithContext(event_api.EventType_PolicyEvaluationAuditEvent, "pdp",
		&event_api.EventContext{
			OrgId:     "0",
			ProjectId: "0",
		})

	var violations []*event_api.PolicyEvaluationEvent_Data_Result_Violation = make([]*event_api.PolicyEvaluationEvent_Data_Result_Violation, 0)
	event := &event_api.PolicyEvaluationEvent{
		Header:    eh,
		Timestamp: time.Now().UnixMilli(),
		Data: &event_api.PolicyEvaluationEvent_Data{
			Artefact: &event_api.Artefact{
				Ecosystem: artefact.OpenSsfEcosystem(),
				Group:     artefact.Group,
				Name:      artefact.Name,
				Version:   artefact.Version,
			},
			Upstream: &event_api.ArtefactUpstream{
				Type: upstream.Type,
				Name: upstream.Name,
			},
			Result: &event_api.PolicyEvaluationEvent_Data_Result{
				PolicyAllowed:      result.Allow,
				EffectiveAllowed:   gw_allowed,
				MonitorMode:        monitor_mode,
				Violations:         violations,
				PackageQueryStatus: &event_api.PolicyEvaluationEvent_Data_Result_PackageMetaQueryStatus{},
			},
			Username:    userId,
			Enrichments: &event_api.PolicyEvaluationEvent_Data_ArtefactEnrichments{},
		},
	}

	for _, v := range pdsResponse.Licenses {
		event.Data.Enrichments.Licenses = append(event.Data.Enrichments.Licenses, v.Id)
	}

	for _, v := range pdsResponse.Vulnerabilities {
		event.Data.Enrichments.Advisories = append(event.Data.Enrichments.Advisories,
			&event_api.PolicyEvaluationEvent_Data_ArtefactEnrichments_ArtefactAdvisory{
				Source:   v.Id.Source,
				SourceId: v.Id.Id,
				Title:    v.Name,
				Severity: v.Severity,
			})
	}

	for _, v := range result.Violations {
		event.Data.Result.Violations = append(event.Data.Result.Violations,
			&event_api.PolicyEvaluationEvent_Data_Result_Violation{
				Code:    int32(v.Code),
				Message: v.Message,
			})
	}

	// status.FromError takes care of handling non-grpc error as well
	grpcStatus, _ := grpc_err_status.FromError(enrichmentErr)
	event.Data.Result.PackageQueryStatus.Code = grpcStatus.Code().String()
	event.Data.Result.PackageQueryStatus.Message = grpcStatus.Message()

	log.Printf("Event: %v", event)

	topic := s.config.Global.PdpService.Publisher.TopicMappings["policy_audit"]
	err := s.messagingService.Publish(topic, event)
	if err != nil {
		log.Printf("[ERROR] Failed to publish audit event to topic: %s err: %v", topic, err)
	}
}
