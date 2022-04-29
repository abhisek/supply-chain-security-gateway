package pds

import (
	"context"
	"errors"
	"log"

	api "github.com/abhisek/supply-chain-gateway/services/gen"
	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db"
)

type policyDataServer struct {
	api.PolicyDataServiceServer
	config     *common_config.Config
	repository *db.VulnerabilityRepository
}

func NewPolicyDataService(config *common_config.Config, repo *db.VulnerabilityRepository) (api.PolicyDataServiceServer, error) {
	return &policyDataServer{
		config:     config,
		repository: repo,
	}, nil
}

func (s *policyDataServer) FindVulnerabilitiesByArtefact(ctx context.Context,
	req *api.FindVulnerabilityByArtefactRequest) (*api.VulnerabilityList, error) {

	log.Printf("Handling query req for: %s/%s", req.Artefact.Ecosystem, req.Artefact.Name)

	vulns := []*api.VulnerabilityMeta{}
	return &api.VulnerabilityList{
		Vulnerabilities: vulns,
	}, nil
}

func (s *policyDataServer) GetVulnerabilityDetails(ctx context.Context,
	req *api.GetVulnerabilityByIdRequest) (*api.VulnerabilityDetail, error) {
	return &api.VulnerabilityDetail{}, errors.New("unimplemented endpoint")
}
