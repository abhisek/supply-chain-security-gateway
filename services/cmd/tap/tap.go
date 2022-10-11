package main

import (
	"context"
	"log"

	common_adapters "github.com/abhisek/supply-chain-gateway/services/pkg/common/adapters"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/logger"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/messaging"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/obs"
	"github.com/abhisek/supply-chain-gateway/services/pkg/tap"

	envoy_v3_ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"

	"google.golang.org/grpc"
)

func main() {
	logger.Init("tap")
	config.Bootstrap("", true)

	tracerShutDown := obs.InitTracing()
	defer tracerShutDown(context.Background())

	messageAdapter, err := config.GetMessagingConfigByName(config.
		TapServiceConfig().PublisherConfig.MessagingAdapterName)
	if err != nil {
		logger.Fatalf("failed to get messaging config: %v", err)
	}

	msgService, err := messaging.NewNatsMessagingService(messageAdapter)
	if err != nil {
		log.Fatalf("Failed to create messaging service: %v", err)
	}

	tapService, err := tap.NewTapService(msgService, []tap.TapHandlerRegistration{
		tap.NewTapEventPublisherRegistration(msgService),
	})

	if err != nil {
		log.Fatalf("Failed to create tap service: %s", err.Error())
	}

	common_adapters.StartGrpcServer("TAP", "0.0.0.0", "9001",
		[]grpc.ServerOption{grpc.MaxConcurrentStreams(5000)}, func(s *grpc.Server) {
			envoy_v3_ext_proc_pb.RegisterExternalProcessorServer(s, tapService)
		})
}
