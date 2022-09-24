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

		if response.ProjectScorecard.Checks != nil {
			pdsResponse.Scorecard.Checks[openssf.ScBinaryArtifactsCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.BinaryArtifacts)
			pdsResponse.Scorecard.Checks[openssf.ScBranchProtectionCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.BranchProtection)
			pdsResponse.Scorecard.Checks[openssf.ScCiiBestPracticeCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.CiiBestPractices)
			pdsResponse.Scorecard.Checks[openssf.ScCodeReviewCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.CodeReview)
			pdsResponse.Scorecard.Checks[openssf.ScDangerousWorkflowCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.DangerousWorkflow)
			pdsResponse.Scorecard.Checks[openssf.ScDependencyUpdateToolCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.DependencyUpdateTool)
			pdsResponse.Scorecard.Checks[openssf.ScFuzzingCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.Fuzzing)
			pdsResponse.Scorecard.Checks[openssf.ScLicenseCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.License)
			pdsResponse.Scorecard.Checks[openssf.ScMaintainedCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.Maintained)
			pdsResponse.Scorecard.Checks[openssf.ScPackagingCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.Packaging)
			pdsResponse.Scorecard.Checks[openssf.ScPinnedDependenciesCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.PinnedDependencies)
			pdsResponse.Scorecard.Checks[openssf.ScSastCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.Sast)
			pdsResponse.Scorecard.Checks[openssf.ScSecurityPolicyCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.SecurityPolicy)
			pdsResponse.Scorecard.Checks[openssf.ScSignedReleasesCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.SignedReleases)
			pdsResponse.Scorecard.Checks[openssf.ScTokenPermissionsCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.TokenPermissions)
			pdsResponse.Scorecard.Checks[openssf.ScVulnerabilitiesCheck] =
				rayaScorecardCheckToOpenSsfScorecardCheck(response.ProjectScorecard.Checks.Vulnerabilities)
		}
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
