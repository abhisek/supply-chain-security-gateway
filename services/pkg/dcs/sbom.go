package dcs

import (
	"log"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
)

const (
	sbomCollectorGroupName = "sbom-collector-group"
	sbomCollectorName      = "SBOM Data Collector"
)

type sbomCollector struct{}

func sbomCollectorSubscription() eventSubscription[common_models.Artefact] {
	h := &sbomCollector{}
	return h.subscription()
}

func (s *sbomCollector) subscription() eventSubscription[common_models.Artefact] {
	cfg, err := config.Current().GC()
	if err != nil {
		panic("service config is not bootstrapped")
	}

	return eventSubscription[common_models.Artefact]{
		name:    sbomCollectorName,
		group:   sbomCollectorGroupName,
		topic:   cfg.Services.Tap.PublisherConfig.TopicNames.UpstreamRequest,
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
