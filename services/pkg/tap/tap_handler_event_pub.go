package tap

import (
	"context"
	"fmt"
	"log"

	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/messaging"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"

	envoy_v3_ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
)

type tapEventPublisher struct {
	config           *common_config.Config
	messagingService messaging.MessagingService
}

func NewTapEventPublisherRegistration(config *common_config.Config, msgService messaging.MessagingService) TapHandlerRegistration {
	return TapHandlerRegistration{
		ContinueOnError: true,
		Handler:         &tapEventPublisher{config: config, messagingService: msgService},
	}
}

func (h *tapEventPublisher) HandleRequestHeaders(ctx context.Context,
	req *envoy_v3_ext_proc_pb.ProcessingRequest_RequestHeaders) error {

	log.Printf("Publishing request headers event")
	path, err := findHeaderValue(req, "path")
	if err != nil {
		log.Printf("Failed to publish: %v", err)
	}

	artefact, err := common_models.GetArtefactByHostAndPath(h.config.Global.Upstreams, "", path)
	if err != nil {
		return fmt.Errorf("Failed to resolve artefact")
	}

	topic := h.config.Global.TapService.Publisher.TopicMappings["upstream_request"]
	event := common_models.NewArtefactRequestEvent(artefact)

	err = h.messagingService.Publish(topic, event)
	if err != nil {
		log.Printf("Error publishing event: %v", err)
		return err
	}

	return nil
}

func (h *tapEventPublisher) HandleResponseHeaders(ctx context.Context,
	req *envoy_v3_ext_proc_pb.ProcessingRequest_ResponseHeaders) error {

	log.Printf("Publishing response headers event - NOP")
	return nil
}
