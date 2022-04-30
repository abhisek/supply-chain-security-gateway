package pdp

import common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"

func NewPolicyInput(target common_models.Artefact,
	upstream common_models.ArtefactUpStream,
	vulnerabilities []common_models.ArtefactVulnerability) PolicyInput {
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
			Vulnerabilities: []PolicyEvalTargetVulnerability{vulnerabilities...},
		},
	}
}

func (s PolicyResponse) Allowed() bool {
	return (s.Allow) && (len(s.Violations) == 0)
}
