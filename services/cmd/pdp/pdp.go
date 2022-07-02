package main

import (
	"log"

	common_adapters "github.com/abhisek/supply-chain-gateway/services/pkg/common/adapters"
	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"

	"github.com/abhisek/supply-chain-gateway/services/pkg/pdp"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"
)

func main() {
	config, err := common_config.LoadGlobal("")
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

	policyDataService, err := pdp.NewPolicyDataServiceClient(config.Global.PdpService)
	if err != nil {
		log.Fatalf("Failed to create policy data service client: %v", err)
	}

	authService, err := pdp.NewAuthorizationService(config, policyDataService)
	if err != nil {
		log.Fatalf("Failed to create auth service: %s", err.Error())
	}

	common_adapters.StartGrpcServer("PDP", "0.0.0.0", "9000",
		[]grpc.ServerOption{}, func(s *grpc.Server) {
			envoy_service_auth_v3.RegisterAuthorizationServer(s, authService)
		})
}
