package pdp

import (
	"context"

	raya_api "github.com/abhisek/supply-chain-gateway/services/gen"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
)

type pdsRayaClient struct {
	client raya_api.RayaClient
}

func (pds *pdsRayaClient) GetPackageMetaByVersion(ctx context.Context,
	ecosystem, group, name, version string) ([]common_models.ArtefactVulnerability, error) {
	return []common_models.ArtefactVulnerability{}, nil
}
