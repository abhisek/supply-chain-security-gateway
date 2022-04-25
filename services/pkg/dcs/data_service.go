package dcs

import (
	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/messaging"
)

type DataCollectionService struct {
	config           *common_config.Config
	messagingService messaging.MessagingService
}

func NewDataCollectionService(config *common_config.Config, msgService messaging.MessagingService) (*DataCollectionService, error) {
	return &DataCollectionService{config: config, messagingService: msgService}, nil
}

func (svc *DataCollectionService) Start() {
	registerSubscriber(svc.messagingService, sbomCollectorSubscription(svc.config))
	waitForSubscribers()
}
