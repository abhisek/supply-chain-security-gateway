package pdp

import (
	"context"

	pds_api "github.com/abhisek/supply-chain-gateway/services/gen"
	raya_api "github.com/abhisek/supply-chain-gateway/services/gen"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
	"google.golang.org/grpc"
)

type PolicyDataClientInterface interface {
	GetPackageMetaByVersion(ctx context.Context, ecosystem, group, name, version string) ([]common_models.ArtefactVulnerability, error)
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
