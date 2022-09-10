package main

import (
	"context"
	"os"

	common_adapters "github.com/abhisek/supply-chain-gateway/services/pkg/common/adapters"
	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/logger"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/messaging"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/obs"

	"github.com/abhisek/supply-chain-gateway/services/pkg/pdp"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"
)

func main() {
	logger.Init("pdp")

	tracerShutDown := obs.InitTracing()
	defer tracerShutDown(context.Background())

	config, err := common_config.LoadGlobal("")
	if err != nil {
		logger.Fatalf("Failed to load config: %s", err.Error())
	}

	policyDataService, err := pdp.NewPolicyDataServiceClient(config.Global.PdpService)
	if err != nil {
		logger.Fatalf("Failed to create policy data service client: %v", err)
	}

	messagingService, err := buildMessagingService(config)
	if err != nil {
		logger.Fatalf("Failed to build messaging service: %v", err)
	}

	authService, err := pdp.NewAuthorizationService(config, policyDataService, messagingService)
	if err != nil {
		logger.Fatalf("Failed to create auth service: %s", err.Error())
	}

	common_adapters.StartGrpcServer("PDP", "0.0.0.0", "9000",
		[]grpc.ServerOption{}, func(s *grpc.Server) {
			envoy_service_auth_v3.RegisterAuthorizationServer(s, authService)
		})
}

func buildMessagingService(config *common_config.Config) (messaging.MessagingService, error) {
	switch config.Global.PdpService.Publisher.Type {
	case "kafka-pongo":
		logger.Infof("Using Kafka (pongo) messaging service")
		return messaging.NewKafkaProtobufMessagingService(os.Getenv("PDP_KAFKA_PONGO_BOOTSTRAP_SERVERS"),
			os.Getenv("PDP_KAFKA_PONGO_SCHEMA_REGISTRY_URL"))
	default:
		logger.Infof("Using NATs messaging service")
		return messaging.NewNatsMessagingService(config)
	}
}
