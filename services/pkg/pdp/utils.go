package pdp

import (
	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
)

// Stop gap method to map a spec based upstream into legacy upstream
func toLegacyUpstream(us *config_api.GatewayUpstream) common_models.ArtefactUpStream {
	return common_models.ToUpstream(us)
}

func isMonitorMode() bool {
	return config.PdpServiceConfig().MonitorMode
}
