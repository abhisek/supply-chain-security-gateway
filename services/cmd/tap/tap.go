package main

import (
	"log"

	common_adapters "github.com/abhisek/supply-chain-gateway/services/pkg/common/adapters"
	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/tap"

	envoy_v3_ext_proc_pb "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"

	"google.golang.org/grpc"
)

func main() {
	config, err := common_config.LoadGlobal("")
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

	tapService, err := tap.NewTapService(config, []tap.TapHandlerRegistration{
		tap.NewTapEventPublisherRegistration(),
	})

	if err != nil {
		log.Fatalf("Failed to create tap service: %s", err.Error())
	}

	common_adapters.StartGrpcServer("TAP", "0.0.0.0", "9001",
		[]grpc.ServerOption{grpc.MaxConcurrentStreams(5000)}, func(s *grpc.Server) {
			envoy_v3_ext_proc_pb.RegisterExternalProcessorServer(s, tapService)
		})
}
