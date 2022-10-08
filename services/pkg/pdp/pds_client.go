package pdp

import (
	"context"
	"fmt"

	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
	pds_api "github.com/abhisek/supply-chain-gateway/services/gen"
	raya_api "github.com/abhisek/supply-chain-gateway/services/gen"
	common_adapters "github.com/abhisek/supply-chain-gateway/services/pkg/common/adapters"
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

func NewPolicyDataServiceClient(cfg *config_api.PdsClientConfig) (PolicyDataClientInterface, error) {
	grpconn, err := buildGrpcClient(cfg.GetCommon().GetHost(),
		fmt.Sprint(cfg.GetCommon().GetPort()), cfg.GetCommon().GetMtls())
	if err != nil {
		return nil, err
	}

	switch cfg.Type {
	case config_api.PdsClientType_LOCAL:
		return NewLocalPolicyDataClient(grpconn), nil
	case config_api.PdsClientType_RAYA:
		return NewRayaPolicyDataServiceClient(grpconn), nil
	default:
		return nil, fmt.Errorf("unknown pds client type:%s", cfg.Type.String())
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
