package pdp

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	raya_api "github.com/abhisek/supply-chain-gateway/services/gen"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
)

type pdsRayaClient struct {
	client raya_api.RayaClient
}

func (pds *pdsRayaClient) GetPackageMetaByVersion(ctx context.Context,
	ecosystem, group, name, version string) ([]common_models.ArtefactVulnerability, error) {
	vulnerabilities := []common_models.ArtefactVulnerability{}
	pkgName := ""

	if !utils.IsEmptyString(group) {
		pkgName = fmt.Sprintf("%s:%s", group, name)
	} else {
		pkgName = name
	}

	request := &raya_api.PackageVersionMetaQueryRequest{
		PackageVersion: &raya_api.PackageVersion{
			Package: &raya_api.Package{
				Ecosystem: rayaEcosystemName(ecosystem),
				Name:      pkgName,
			},
			Version: version,
		},
	}

	log.Printf("Querying Raya with: %v", request)

	response, err := pds.client.GetPackageMetaByVersion(ctx, request)
	if err != nil {
		return vulnerabilities, err
	}

	severityMapper := func(s raya_api.Severity) string {
		switch s {
		case raya_api.Severity_CRITICAL:
			return common_models.ArtefactVulnerabilitySeverityCritical
		case raya_api.Severity_HIGH:
			return common_models.ArtefactVulnerabilitySeverityHigh
		case raya_api.Severity_MEDIUM:
			return common_models.ArtefactVulnerabilitySeverityMedium
		case raya_api.Severity_LOW:
			return common_models.ArtefactVulnerabilitySeverityLow
		default:
			return common_models.ArtefactVulnerabilitySeverityInfo
		}
	}

	for _, adv := range response.Advisories {
		if adv == nil {
			continue
		}

		vulnerabilities = append(vulnerabilities, common_models.ArtefactVulnerability{
			Id: common_models.ArtefactVulnerabilityId{
				Source: adv.Source,
				Id:     adv.SourceId,
			},
			Name:     adv.Title,
			Severity: severityMapper(adv.AdvisorySeverity.Severity),
			Scores: []common_models.ArtefactVulnerabilityScore{
				{
					Type:  common_models.ArtefactVulnerabilityScoreTypeCVSSv3,
					Value: strconv.FormatFloat(float64(adv.AdvisorySeverity.Cvssv3Score), 'f', -1, 32),
				},
			},
		})
	}

	return vulnerabilities, nil
}

func rayaEcosystemName(name string) string {
	return strings.ToUpper(name)
}
