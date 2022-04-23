package tap

import (
	"context"
	"fmt"
	"log"
	"strings"

	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"

	envoy_v3_ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
)

type tapEventPublisher struct {
	config *common_config.Config
}

func NewTapEventPublisherRegistration(config *common_config.Config) TapHandlerRegistration {
	return TapHandlerRegistration{
		ContinueOnError: true,
		Handler:         &tapEventPublisher{config: config},
	}
}

func (h *tapEventPublisher) HandleRequestHeaders(ctx context.Context,
	req *envoy_v3_ext_proc_pb.ProcessingRequest_RequestHeaders,
	resp *envoy_v3_ext_proc_pb.ProcessingResponse) error {

	log.Printf("Publishing request headers event")
	path, err := findHeaderValue(req, "path")
	if err != nil {
		log.Printf("Failed to publish: %v", err)
	}

	artefact, err := common_models.GetArtefactByHostAndPath(h.config.Global.Upstreams, "", path)
	if err != nil {
		return fmt.Errorf("Failed to resolve artefact")
	}

	event := common_models.NewArtefactRequestEvent(artefact)
	fmt.Printf("event: %v\n", event)

	return nil
}

func (h *tapEventPublisher) HandleResponseHeaders(ctx context.Context,
	req *envoy_v3_ext_proc_pb.ProcessingRequest_ResponseHeaders,
	resp *envoy_v3_ext_proc_pb.ProcessingResponse) error {

	log.Printf("Publishing response headers event - NOP")
	return nil
}

// Header keys are stored as ":key" by envoy
func findHeaderValue(req *envoy_v3_ext_proc_pb.ProcessingRequest_RequestHeaders,
	key string) (string, error) {
	for _, h := range req.RequestHeaders.Headers.Headers {
		if strings.EqualFold(":"+key, h.Key) {
			return h.Value, nil
		}
	}

	return "", fmt.Errorf("header with key: %s not found", key)
}
