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

	if response.ProjectScorecard != nil {
		pdsResponse.Scorecard.Timestamp = response.ProjectScorecard.Timestamp
		pdsResponse.Scorecard.Score = response.ProjectScorecard.Score
		pdsResponse.Scorecard.Version = response.ProjectScorecard.Version

		if response.ProjectScorecard.Repo != nil {
			pdsResponse.Scorecard.Repo.Name = response.ProjectScorecard.Repo.Name
			pdsResponse.Scorecard.Repo.Commit = response.ProjectScorecard.Repo.Commit
		}

		checksMap := map[string]openssf.ProjectScorecardCheck{}
		if response.ProjectScorecard.Checks != nil {
			checksMap[openssf.ScBinaryArtifactsCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.BinaryArtifacts)
			checksMap[openssf.ScBranchProtectionCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.BranchProtection)
			checksMap[openssf.ScCiiBestPracticeCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.CiiBestPractices)
			checksMap[openssf.ScCodeReviewCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.CodeReview)
			checksMap[openssf.ScDangerousWorkflowCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.DangerousWorkflow)
			checksMap[openssf.ScDependencyUpdateToolCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.DependencyUpdateTool)
			checksMap[openssf.ScFuzzingCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.Fuzzing)
			checksMap[openssf.ScLicenseCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.License)
			checksMap[openssf.ScMaintainedCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.Maintained)
			checksMap[openssf.ScPackagingCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.Packaging)
			checksMap[openssf.ScPinnedDependenciesCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.PinnedDependencies)
			checksMap[openssf.ScSastCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.Sast)
			checksMap[openssf.ScSecurityPolicyCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.SecurityPolicy)
			checksMap[openssf.ScSignedReleasesCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.SignedReleases)
			checksMap[openssf.ScTokenPermissionsCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.TokenPermissions)
			checksMap[openssf.ScVulnerabilitiesCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.Vulnerabilities)
		}

		pdsResponse.Scorecard.Checks = checksMap
	}

	return pdsResponse, nil
}

func rayaScorecardCheckToOpenSsfScorecardCheck(sc *raya_api.ProjectScorecardCheck) openssf.ProjectScorecardCheck {
	psc := openssf.ProjectScorecardCheck{}
	if sc == nil {
		return psc
	}

	psc.Reason = sc.Reason
	psc.Score = sc.Score

	return psc
}

func rayaEcosystemName(name string) string {
	return strings.ToUpper(name)
}
