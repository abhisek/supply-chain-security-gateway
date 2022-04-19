package tap

import (
	"context"
	"log"

	envoy_v3_ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
)

type tapEventPublisher struct{}

func NewTapEventPublisherRegistration() TapHandlerRegistration {
	return TapHandlerRegistration{
		ContinueOnError: true,
		Handler:         &tapEventPublisher{},
	}
}

func (h *tapEventPublisher) HandleRequestHeaders(ctx context.Context,
	req *envoy_v3_ext_proc_pb.ProcessingRequest_RequestHeaders,
	resp *envoy_v3_ext_proc_pb.ProcessingResponse) error {

	log.Printf("Publishing request headers event")
	return nil
}

func (h *tapEventPublisher) HandleResponseHeaders(ctx context.Context,
	req *envoy_v3_ext_proc_pb.ProcessingRequest_ResponseHeaders,
	resp *envoy_v3_ext_proc_pb.ProcessingResponse) error {

	log.Printf("Publishing response headers event")
	return nil
}
