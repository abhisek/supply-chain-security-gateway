package pdp

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	raya_api "github.com/abhisek/supply-chain-gateway/services/gen"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/openssf"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
)

type pdsRayaClient struct {
	client raya_api.RayaClient
}

func (pds *pdsRayaClient) GetPackageMetaByVersion(ctx context.Context,
	ecosystem, group, name, version string) (PolicyDataServiceResponse, error) {
	pkgName := ""

	if !utils.IsEmptyString(group) && ecosystem == openssf.VulnerabilityEcosystemMaven {
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
		return PolicyDataServiceResponse{}, err
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

	pdsResponse := PolicyDataServiceResponse{}
	for _, adv := range response.Advisories {
		if adv == nil {
			continue
		}

		pdsResponse.Vulnerabilities = append(pdsResponse.Vulnerabilities, common_models.ArtefactVulnerability{
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

	for _, license := range response.Licenses {
		pdsResponse.Licenses = append(pdsResponse.Licenses, common_models.ArtefactLicense{
			Type: common_models.ArtefactLicenseTypeSpdx,
			Id:   license,
			Name: license,
		})
	}

	return pdsResponse, nil
}

func rayaEcosystemName(name string) string {
	return strings.ToUpper(name)
}
