package main

import (
	"log"
	"os"

	common_adapters "github.com/abhisek/supply-chain-gateway/services/pkg/common/adapters"
	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"

	pds_api "github.com/abhisek/supply-chain-gateway/services/gen"

	"github.com/abhisek/supply-chain-gateway/services/pkg/pdp"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"
)

func main() {
	config, err := common_config.LoadGlobal("")
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

	grpconn, err := common_adapters.GrpcMtlsClient("PDS", os.Getenv("PDS_HOST"),
		os.Getenv("PDS_HOST"), os.Getenv("PDS_PORT"), []grpc.DialOption{}, func(conn *grpc.ClientConn) {})
	if err != nil {
		log.Fatalf("Failed to establish connection with PDS: %v", err)
	}

	policyDataService := pds_api.NewPolicyDataServiceClient(grpconn)
	authService, err := pdp.NewAuthorizationService(config, policyDataService)
	if err != nil {
		log.Fatalf("Failed to create auth service: %s", err.Error())
	}

	common_adapters.StartGrpcServer("PDP", "0.0.0.0", "9000",
		[]grpc.ServerOption{}, func(s *grpc.Server) {
			envoy_service_auth_v3.RegisterAuthorizationServer(s, authService)
		})
}
