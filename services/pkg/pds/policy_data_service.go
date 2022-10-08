package pds

import (
	"context"
	"errors"
	"log"

	api "github.com/abhisek/supply-chain-gateway/services/gen"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type policyDataServer struct {
	api.PolicyDataServiceServer
	repository *db.VulnerabilityRepository
}

func NewPolicyDataService(repo *db.VulnerabilityRepository) (api.PolicyDataServiceServer, error) {
	return &policyDataServer{
		repository: repo,
	}, nil
}

func (s *policyDataServer) FindVulnerabilitiesByArtefact(ctx context.Context,
	req *api.FindVulnerabilityByArtefactRequest) (*api.VulnerabilityList, error) {

	vulnList := &api.VulnerabilityList{
		Vulnerabilities: []*api.VulnerabilityMeta{},
	}

	log.Printf("Handling query req for: %s/%s", req.Artefact.Ecosystem, req.Artefact.Name)

	// Lookup all vulnerabilities by Ecosystem, group, name
	dbVulnerabilities, err := s.repository.Lookup(req.Artefact.Ecosystem, req.Artefact.Group,
		req.Artefact.Name)
	if err != nil {
		return vulnList, status.Errorf(codes.Internal, "failed on query: %v", err)
	}

	// Fuzzy match to find vulnerabilities applicable for the version
	for _, dbVuln := range dbVulnerabilities {
		wrappedVuln, err := wrapVuln(dbVuln)
		if err != nil {
			log.Printf("failed to wrap db vuln: %v", err)
			continue
		}

		if s.match(wrappedVuln, req.Artefact) {
			vulnList.Vulnerabilities = append(vulnList.Vulnerabilities, &api.VulnerabilityMeta{
				Id:       dbVuln.ExternalId,
				Source:   dbVuln.ExternalSource,
				Title:    dbVuln.Title,
				Severity: wrappedVuln.FriendlySeverityCode(),
				Scores:   wrappedVuln.Severity(),
			})
		}
	}

	log.Printf("Supplying vulnerabilities: %s", utils.Introspect(vulnList))
	return vulnList, nil
}

func (s *policyDataServer) GetVulnerabilityDetails(ctx context.Context,
	req *api.GetVulnerabilityByIdRequest) (*api.VulnerabilityDetail, error) {
	return &api.VulnerabilityDetail{}, errors.New("unimplemented endpoint")
}

// TODO: Fuzzy match
func (s *policyDataServer) match(vuln VulnerabilitySchemaWrapper, artefact *api.Artefact) bool {
	affects := vuln.Affects()
	for _, a := range affects {
		for _, v := range a.Versions {
			if artefact.Version == v {
				return true
			}
		}
	}

	return false
}
