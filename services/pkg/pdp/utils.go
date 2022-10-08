package pdp

import (
	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
)

// Stop gap method to map a spec based upstream into legacy upstream
func toLegacyUpstream(us *config_api.GatewayUpstream) common_models.ArtefactUpStream {
	upstream := common_models.ArtefactUpStream{
		Name: us.Name,
		Type: us.Type.String(),
		RoutingRule: common_models.ArtefactRoutingRule{
			Prefix: us.Route.PathPrefix,
			Host:   us.Route.Host,
		},
		Authentication: common_models.ArtefactUpstreamAuthentication{
			Type:     us.Repository.Authentication.Type.String(),
			Provider: us.Repository.Authentication.Provider,
		},
	}

	return upstream
}
