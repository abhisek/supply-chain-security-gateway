package pdp

import (
	"github.com/abhisek/supply-chain-gateway/services/pkg/auth"
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
)

func NewPolicyInput(target common_models.Artefact,
	upstream common_models.ArtefactUpStream,
	requester auth.AuthenticatedIdentity,
	enrichments PolicyDataServiceResponse) PolicyInput {

	vulns := []PolicyEvalTargetVulnerability{}
	for _, v := range enrichments.Vulnerabilities {
		vulns = append(vulns, PolicyEvalTargetVulnerability{v})
	}

	lics := []PolicyEvalTargetLicense{}
	for _, l := range enrichments.Licenses {
		lics = append(lics, PolicyEvalTargetLicense{l})
	}

	return PolicyInput{
		Kind: policyInputKind,
		Version: PolicyInputVersion{
			Major: policyInputMajorVersion,
			Minor: policyInputMinorVersion,
			Patch: policyInputPatchVersion,
		},
		Target: PolicyInputTarget{
			Artefact:        PolicyEvalTargetArtefact{target},
			Upstream:        PolicyEvalTargetUpstream{upstream},
			Vulnerabilities: vulns,
			Licenses:        lics,
		},
		Principal: PolicyInputPrincipal{
			UserId:    requester.UserId(),
			ProjectId: requester.ProjectId(),
			OrgId:     requester.OrgId(),
		},
	}
}

func (s PolicyResponse) Allowed() bool {
	return (s.Allow) && (len(s.Violations) == 0)
}
