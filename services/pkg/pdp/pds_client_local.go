package pdp

import (
	"context"

	pds_api "github.com/abhisek/supply-chain-gateway/services/gen"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
)

type pdsLocalImplementation struct {
	client pds_api.PolicyDataServiceClient
}

func (pds *pdsLocalImplementation) GetPackageMetaByVersion(ctx context.Context,
	ecosystem, group, name, version string) (PolicyDataServiceResponse, error) {
	resp, err := pds.client.FindVulnerabilitiesByArtefact(ctx, &pds_api.FindVulnerabilityByArtefactRequest{
		Artefact: &pds_api.Artefact{
			Ecosystem: ecosystem,
			Group:     group,
			Name:      name,
			Version:   version,
		},
	})

	if err != nil {
		return PolicyDataServiceResponse{}, err
	}

	return PolicyDataServiceResponse{
		Vulnerabilities: pds.mapVulnerabilities(resp.Vulnerabilities),
	}, nil
}

func (pds *pdsLocalImplementation) mapVulnerabilities(src []*pds_api.VulnerabilityMeta) []common_models.ArtefactVulnerability {
	target := []common_models.ArtefactVulnerability{}

	if src == nil || len(src) == 0 {
		return target
	}

	for _, s := range src {
		mv := common_models.ArtefactVulnerability{
			Name: s.Title,
			Id: common_models.ArtefactVulnerabilityId{
				Source: s.Source,
				Id:     s.Id,
			},
			Scores: []common_models.ArtefactVulnerabilityScore{},
		}

		switch s.Severity {
		case pds_api.VulnerabilitySeverity_CRITICAL:
			mv.Severity = common_models.ArtefactVulnerabilitySeverityCritical
			break
		case pds_api.VulnerabilitySeverity_HIGH:
			mv.Severity = common_models.ArtefactVulnerabilitySeverityHigh
			break
		case pds_api.VulnerabilitySeverity_MEDIUM:
			mv.Severity = common_models.ArtefactVulnerabilitySeverityMedium
			break
		case pds_api.VulnerabilitySeverity_LOW:
			mv.Severity = common_models.ArtefactVulnerabilitySeverityLow
		default:
			mv.Severity = common_models.ArtefactVulnerabilitySeverityInfo
		}

		for _, score := range s.Scores {
			mv.Scores = append(mv.Scores, common_models.ArtefactVulnerabilityScore{
				Type:  score.Type,
				Value: score.Value,
			})
		}

		target = append(target, mv)
	}

	return target
}
