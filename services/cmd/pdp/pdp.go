package main

import (
	"context"
	"fmt"
	"os"

	common_adapters "github.com/abhisek/supply-chain-gateway/services/pkg/common/adapters"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/logger"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/messaging"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/obs"

	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
	"github.com/abhisek/supply-chain-gateway/services/pkg/pdp"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"
)

func main() {
	logger.Init("pdp")
	config.Bootstrap("", true)

	tracerShutDown := obs.InitTracing()
	defer tracerShutDown(context.Background())

	policyDataService, err := pdp.NewPolicyDataServiceClient(config.PdpServiceConfig().GetPdsClient())
	if err != nil {
		logger.Fatalf("Failed to create policy data service client: %v", err)
	}

	messagingService, err := buildMessagingService()
	if err != nil {
		logger.Fatalf("Failed to build messaging service: %v", err)
	}

	authService, err := pdp.NewAuthorizationService(policyDataService, messagingService)
	if err != nil {
		logger.Fatalf("Failed to create auth service: %s", err.Error())
	}

	common_adapters.StartGrpcServer("PDP", "0.0.0.0", "9000",
		[]grpc.ServerOption{}, func(s *grpc.Server) {
			envoy_service_auth_v3.RegisterAuthorizationServer(s, authService)
		})
}

func buildMessagingService() (messaging.MessagingService, error) {
	cfg := config.PdpServiceConfig()
	messageAdapter, err := config.GetMessagingConfigByName(cfg.PublisherConfig.MessagingAdapterName)
	if err != nil {
		return nil, err
	}

	// FIXME: Migrate to messaging.NewService(...) config based factory
	switch messageAdapter.Type {
	case config_api.MessagingAdapter_KAFKA:
		logger.Infof("Using Kafka (pongo) messaging service")
		return messaging.NewKafkaProtobufMessagingService(os.Getenv("PDP_KAFKA_PONGO_BOOTSTRAP_SERVERS"),
			os.Getenv("PDP_KAFKA_PONGO_SCHEMA_REGISTRY_URL"))
	case config_api.MessagingAdapter_NATS:
		logger.Infof("Using NATs messaging service")
		return messaging.NewNatsMessagingService(messageAdapter)
	default:
		return nil, fmt.Errorf("unknown message adapter type: %s", messageAdapter.Type.String())
	}
}
