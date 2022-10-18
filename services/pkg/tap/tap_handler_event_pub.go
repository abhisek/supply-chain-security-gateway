package tap

import (
	"context"
	"fmt"
	"log"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/messaging"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"

	envoy_v3_ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
)

type tapEventPublisher struct {
	messagingService messaging.MessagingService
}

func NewTapEventPublisherRegistration(msgService messaging.MessagingService) TapHandlerRegistration {
	return TapHandlerRegistration{
		ContinueOnError: true,
		Handler:         &tapEventPublisher{messagingService: msgService},
	}
}

func (h *tapEventPublisher) HandleRequestHeaders(ctx context.Context,
	req *envoy_v3_ext_proc_pb.ProcessingRequest_RequestHeaders) error {

	cfg := config.TapServiceConfig()

	log.Printf("Publishing request headers event")

	host, path, err := findHostAndPath(req)
	if err != nil {
		return err
	}

	artefact, err := common_models.GetArtefactByHostAndPath(host, path)
	if err != nil {
		return fmt.Errorf("Failed to resolve artefact")
	}

	topic := cfg.GetPublisherConfig().GetTopicNames().GetUpstreamRequest()

	// TODO: Migrate this to event spec
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
