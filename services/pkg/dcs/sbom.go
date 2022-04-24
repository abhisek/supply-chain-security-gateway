package dcs

import (
	"log"

	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
)

const (
	sbomCollectorGroupName = "sbom-collector-group"
)

func sbomCollectorSubscription(config *common_config.Config) eventSubscription[common_models.Artefact] {
	return eventSubscription[common_models.Artefact]{
		name:    "SBOM Data Collector",
		topic:   config.Global.TapService.Publisher.TopicMappings["upstream_request"],
		group:   sbomCollectorGroupName,
		handler: sbomCollectorHandler(),
	}
}

func sbomCollectorHandler() eventSubscriptionHandler[common_models.Artefact] {
	return func(event common_models.DomainEvent[common_models.Artefact]) error {
		log.Printf("SBOM collector - Handling artefact: %v", event.Data)
		return nil
	}
}
