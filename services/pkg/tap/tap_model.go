package tap

import (
	"context"

	envoy_v3_ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
)

type TapHandler interface {
	HandleRequestHeaders(context.Context,
		*envoy_v3_ext_proc_pb.ProcessingRequest_RequestHeaders) error
	HandleResponseHeaders(context.Context,
		*envoy_v3_ext_proc_pb.ProcessingRequest_ResponseHeaders) error
}

type TapHandlerRegistration struct {
	ContinueOnError bool
	Handler         TapHandler
}

type TapHandlerChain struct {
	Handlers []TapHandlerRegistration
}
