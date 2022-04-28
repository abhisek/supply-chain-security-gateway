package dcs

import (
	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/messaging"
)

type DataCollectionService struct {
	messagingService        messaging.MessagingService
	config                  *common_config.Config
	vulnerabilityRepository *db.VulnerabilityRepository
}

func NewDataCollectionService(config *common_config.Config,
	msgService messaging.MessagingService,
	vRepo *db.VulnerabilityRepository) (*DataCollectionService, error) {

	return &DataCollectionService{config: config,
		messagingService:        msgService,
		vulnerabilityRepository: vRepo}, nil
}

func (svc *DataCollectionService) Start() {
	registerSubscriber(svc.messagingService, sbomCollectorSubscription(svc.config))
	registerSubscriber(svc.messagingService, vulnCollectorSubscription(svc.config, svc.vulnerabilityRepository))

	waitForSubscribers()
}
