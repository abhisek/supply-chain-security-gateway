package main

import (
	"log"
	"os"
	"strconv"

	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db/adapters"
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

	mysqlPort, err := strconv.ParseInt(os.Getenv("MYSQL_SERVER_PORT"), 0, 16)
	if err != nil {
		log.Fatalf("Failed to parse mysql server port: %v", err)
	}

	mysqlAdapter, err := adapters.NewMySqlAdapter(adapters.MySqlAdapterConfig{
		Host:     os.Getenv("MYSQL_SERVER_HOST"),
		Port:     int16(mysqlPort),
		Username: os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DATABASE"),
	})
	if err != nil {
		log.Fatalf("Failed to initialize MySQL adapter: %v", err)
	}

	err = db.MigrateSqlModels(mysqlAdapter)
	if err != nil {
		log.Fatalf("Failed to run MySQL migration: %v", err)
	}

	repository, err := db.NewVulnerabilityRepository(config, mysqlAdapter)
	if err != nil {
		log.Fatalf("Failed to create vulnerability repository")
	}

	dcs, err := dcs.NewDataCollectionService(config, msgService, repository)
	if err != nil {
		log.Fatalf("Failed to created DCS: %v", err)
	}

	log.Printf("Starting data collector service(s)")
	dcs.Start()
}
