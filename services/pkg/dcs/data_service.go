package dcs

import (
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/messaging"
)

type DataCollectionService struct {
	messagingService        messaging.MessagingService
	vulnerabilityRepository *db.VulnerabilityRepository
}

func NewDataCollectionService(msgService messaging.MessagingService,
	vRepo *db.VulnerabilityRepository) (*DataCollectionService, error) {

	return &DataCollectionService{messagingService: msgService,
		vulnerabilityRepository: vRepo}, nil
}

func (svc *DataCollectionService) Start() {
	registerSubscriber(svc.messagingService, sbomCollectorSubscription())
	registerSubscriber(svc.messagingService, vulnCollectorSubscription(svc.vulnerabilityRepository))

	waitForSubscribers()
}
