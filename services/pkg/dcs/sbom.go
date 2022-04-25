package dcs

import (
	"log"

	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
)

const (
	sbomCollectorGroupName = "sbom-collector-group"
	sbomCollectorName      = "SBOM Data Collector"
)

type sbomCollector struct {
	config *common_config.Config
}

func sbomCollectorSubscription(config *common_config.Config) eventSubscription[common_models.Artefact] {
	h := &sbomCollector{config: config}
	return h.subscription()
}

func (s *sbomCollector) subscription() eventSubscription[common_models.Artefact] {
	return eventSubscription[common_models.Artefact]{
		name:    sbomCollectorName,
		group:   sbomCollectorGroupName,
		topic:   s.config.Global.TapService.Publisher.TopicMappings["upstream_request"],
		handler: s.handler(),
	}
}

func (s *sbomCollector) handler() eventSubscriptionHandler[common_models.Artefact] {
	return func(event common_models.DomainEvent[common_models.Artefact]) error {
		return s.handle(event)
	}
}

func (s *sbomCollector) handle(event common_models.DomainEvent[common_models.Artefact]) error {
	log.Printf("SBOM collector - Handling artefact: %v", event.Data)
	return nil
}
