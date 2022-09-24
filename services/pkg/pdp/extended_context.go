package pdp

import (
	"context"
	"strings"
	"time"

	"github.com/abhisek/supply-chain-gateway/services/pkg/auth"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"go.uber.org/zap"
)

const (
	projectIdHeaderName     = "X-SGW-Project-Id"
	projectEnvHeaderName    = "X-SGW-Project-Env"
	projectLabelsHeaderName = "X-SGW-Project-Labels"
)

type extendedContext struct {
	innerCtx          context.Context
	envoyCheckRequest *envoy_service_auth_v3.CheckRequest
	logger            *zap.SugaredLogger
	identity          auth.AuthenticatedIdentity
	upstream          common_models.ArtefactUpStream
	artefact          common_models.Artefact
}

func ExtendContext(ctx context.Context) *extendedContext {
	return &extendedContext{innerCtx: ctx}
}

func (ctx *extendedContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.innerCtx.Deadline()
}

func (ctx *extendedContext) Done() <-chan struct{} {
	return ctx.innerCtx.Done()
}

func (ctx *extendedContext) Err() error {
	return ctx.innerCtx.Err()
}

func (ctx *extendedContext) Value(key any) any {
	return ctx.innerCtx.Value(key)
}

func (ctx *extendedContext) WithEnvoyCheckRequest(r *envoy_service_auth_v3.CheckRequest) *extendedContext {
	ctx.envoyCheckRequest = r
	return ctx
}

func (ctx *extendedContext) WithLogger(l *zap.SugaredLogger) *extendedContext {
	ctx.logger = l
	return ctx
}

func (ctx *extendedContext) WithAuthIdentity(id auth.AuthenticatedIdentity) *extendedContext {
	ctx.identity = id
	return ctx
}

func (ctx *extendedContext) WithArtefact(artefact common_models.Artefact) *extendedContext {
	ctx.artefact = artefact
	return ctx
}

func (ctx *extendedContext) WithUpstream(upstream common_models.ArtefactUpStream) *extendedContext {
	ctx.upstream = upstream
	return ctx
}

// https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_conn_man/headers
func (ctx *extendedContext) RequestHost() string {
	return ctx.envoyCheckRequest.Attributes.Request.Http.Headers[":authority"]
}

func (ctx *extendedContext) RequestHeader(name string) string {
	return ctx.envoyCheckRequest.Attributes.Request.Http.Headers[":"+strings.ToLower(name)]
}

func (ctx *extendedContext) RequestPath() string {
	return ctx.envoyCheckRequest.Attributes.Request.Http.Path
}

func (ctx *extendedContext) ValidHost() bool {
	parts := strings.SplitN(ctx.RequestHost(), ".", 2)
	return len(parts) == 2
}

func (ctx *extendedContext) GatewayDomain() string {
	return strings.SplitN(ctx.RequestHost(), ".", 2)[0]
}

func (ctx *extendedContext) EnvironmentDomain() string {
	if ctx.ValidHost() {
		return strings.SplitN(ctx.RequestHost(), ".", 2)[1]
	} else {
		return ""
	}
}

func (ctx *extendedContext) UserId() string {
	if ctx.identity == nil {
		return ""
	}

	return ctx.identity.UserId()
}

func (ctx *extendedContext) OrgId() string {
	if ctx.identity == nil {
		return ""
	}

	return ctx.identity.OrgId()
}

// There is a bug here. ProjectId is overridden through env
// but policy engine still sees the projectId encoded in username
func (ctx *extendedContext) ProjectId() string {
	projectId := ctx.RequestHeader(projectIdHeaderName)
	if utils.IsEmptyString(projectId) && (ctx.identity != nil) {
		projectId = ctx.identity.ProjectId()
	}

	return projectId
}

func (ctx *extendedContext) Artefact() common_models.Artefact {
	return ctx.artefact
}

func (ctx *extendedContext) Upstream() common_models.ArtefactUpStream {
	return ctx.upstream
}
