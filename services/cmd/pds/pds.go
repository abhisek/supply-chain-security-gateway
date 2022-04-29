package main

import (
	"log"
	"os"
	"strconv"

	api "github.com/abhisek/supply-chain-gateway/services/gen"

	common_adapters "github.com/abhisek/supply-chain-gateway/services/pkg/common/adapters"
	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"google.golang.org/grpc"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db/adapters"
	"github.com/abhisek/supply-chain-gateway/services/pkg/pds"
)

func main() {
	config, err := common_config.LoadGlobal("")
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
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

	repository, err := db.NewVulnerabilityRepository(config, mysqlAdapter)
	if err != nil {
		log.Fatalf("Failed to create vulnerability repository")
	}

	pdService, err := pds.NewPolicyDataService(config, repository)
	if err != nil {
		log.Fatalf("Failed to create policy data service")
	}

	common_adapters.StartGrpcMtlsServer("PDS", os.Getenv("PDS_SERVER_NAME"), "0.0.0.0", "9002",
		[]grpc.ServerOption{grpc.MaxConcurrentStreams(5000),
			common_adapters.GrpcStreamValidatorInterceptor(),
			common_adapters.GrpcUnaryValidatorInterceptor()}, func(s *grpc.Server) {
			api.RegisterPolicyDataServiceServer(s, pdService)
		})
}
