package main

import (
	"log"

	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/messaging"
	"github.com/abhisek/supply-chain-gateway/services/pkg/dcs"
)

func main() {
	config, err := common_config.LoadGlobal("")
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

	msgService, err := messaging.NewNatsMessagingService(config)
	if err != nil {
		log.Fatalf("Failed to create messaging service: %v", err)
	}

	dcs, err := dcs.NewDataCollectionService(config, msgService)
	if err != nil {
		log.Fatalf("Failed to created DCS: %v", err)
	}

	log.Printf("Starting data collector service(s)")
	dcs.Start()
}
