package pdp

import (
	"context"
	"fmt"

	pds_api "github.com/abhisek/supply-chain-gateway/services/gen"
	raya_api "github.com/abhisek/supply-chain-gateway/services/gen"
	common_adapters "github.com/abhisek/supply-chain-gateway/services/pkg/common/adapters"
	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/openssf"

	"google.golang.org/grpc"
)

const (
	pdsClientTypeLocal = "local"
	pdsClientTypeRaya  = "raya"
)

type PolicyDataServiceResponse struct {
	Vulnerabilities []common_models.ArtefactVulnerability `json:"vulnerabilities"`
	Licenses        []common_models.ArtefactLicense       `json:"licenses"`
	Scorecard       openssf.ProjectScorecard
}

type PolicyDataClientInterface interface {
	GetPackageMetaByVersion(ctx context.Context, ecosystem, group, name, version string) (PolicyDataServiceResponse, error)
}

func NewPolicyDataServiceClient(cfg common_config.PdpServiceConfig) (PolicyDataClientInterface, error) {
	grpconn, err := buildGrpcClient(cfg.PdsClient.Host, cfg.PdsClient.Port, cfg.PdsClient.UseMtls)
	if err != nil {
		return nil, err
	}

	switch cfg.PdsClient.Type {
	case pdsClientTypeLocal:
		return NewLocalPolicyDataClient(grpconn), nil
	case pdsClientTypeRaya:
		return NewRayaPolicyDataServiceClient(grpconn), nil
	default:
		return nil, fmt.Errorf("unknown pds client type:%s", cfg.PdsClient.Type)
	}
}

func buildGrpcClient(host string, port string, mtls bool) (*grpc.ClientConn, error) {
	if mtls {
		return common_adapters.GrpcMtlsClient("pds_secure_client", host, host, port,
			common_adapters.NoGrpcDialOptions, common_adapters.NoGrpcConfigurer)
	} else {
		return common_adapters.GrpcInsecureClient("pds_insecure_client", host, port,
			common_adapters.NoGrpcDialOptions, common_adapters.NoGrpcConfigurer)
	}
}

func NewLocalPolicyDataClient(cc grpc.ClientConnInterface) PolicyDataClientInterface {
	return &pdsLocalImplementation{
		client: pds_api.NewPolicyDataServiceClient(cc),
	}
}

func NewRayaPolicyDataServiceClient(conn grpc.ClientConnInterface) PolicyDataClientInterface {
	return &pdsRayaClient{
		client: raya_api.NewRayaClient(conn),
	}
}
