package main

import (
	"context"
	"os"
	"strconv"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db/adapters"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/logger"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/messaging"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/obs"
	"github.com/abhisek/supply-chain-gateway/services/pkg/dcs"
)

func main() {
	logger.Init("dcs")
	config.Bootstrap("", true)

	tracerShutDown := obs.InitTracing()
	defer tracerShutDown(context.Background())

	msgAdapter, err := config.Current().
		GetMessagingConfigByName(config.Current().
			DcsServiceConfig().GetMessagingAdapterName())
	if err != nil {
		logger.Fatalf("Failed to get messaging adapter config")
	}

	msgService, err := messaging.NewService(msgAdapter)
	if err != nil {
		logger.Fatalf("Failed to create messaging service: %v", err)
	}

	mysqlPort, err := strconv.ParseInt(os.Getenv("MYSQL_SERVER_PORT"), 0, 16)
	if err != nil {
		logger.Fatalf("Failed to parse mysql server port: %v", err)
	}

	mysqlAdapter, err := adapters.NewMySqlAdapter(adapters.MySqlAdapterConfig{
		Host:     os.Getenv("MYSQL_SERVER_HOST"),
		Port:     int16(mysqlPort),
		Username: os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Database: os.Getenv("MYSQL_DATABASE"),
	})
	if err != nil {
		logger.Fatalf("Failed to initialize MySQL adapter: %v", err)
	}

	err = db.MigrateSqlModels(mysqlAdapter)
	if err != nil {
		logger.Fatalf("Failed to run MySQL migration: %v", err)
	}

	repository, err := db.NewVulnerabilityRepository(mysqlAdapter)
	if err != nil {
		logger.Fatalf("Failed to create vulnerability repository")
	}

	dcs, err := dcs.NewDataCollectionService(msgService, repository)
	if err != nil {
		logger.Fatalf("Failed to created DCS: %v", err)
	}

	logger.Infof("Starting data collector service(s)")
	dcs.Start()
}
