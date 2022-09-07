package pdp

import (
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
)

const (
	policyInputKind                 = "PolicyInput"
	policyInputMajorVersion         = 1
	policyInputMinorVersion         = 0
	policyInputPatchVersion         = 0
	policyEvalChangeMonitorInterval = "5s"
)

type PolicyEvalTargetArtefact struct {
	common_models.Artefact
}

type PolicyEvalTargetUpstream struct {
	common_models.ArtefactUpStream
}

type PolicyEvalTargetVulnerability struct {
	common_models.ArtefactVulnerability
}

type PolicyEvalTargetLicense struct {
	common_models.ArtefactLicense
}

type PolicyInputVersion struct {
	Major int8 `json:"major"`
	Minor int8 `json:"minor"`
	Patch int8 `json:"patch"`
}

type PolicyInputPrincipal struct {
	Id string `json:"id"`
}

type PolicyInputTarget struct {
	Artefact        PolicyEvalTargetArtefact        `json:"artefact"`
	Upstream        PolicyEvalTargetUpstream        `json:"upstream"`
	Vulnerabilities []PolicyEvalTargetVulnerability `json:"vulnerabilities"`
	Licenses        []PolicyEvalTargetLicense       `json:"licenses"`
}

type PolicyInput struct {
	Kind      string               `json:"kind"`
	Version   PolicyInputVersion   `json:"version"`
	Target    PolicyInputTarget    `json:"target"`
	Principal PolicyInputPrincipal `json:"principal"`
}

type PolicyViolation struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type PolicyResponse struct {
	Allow      bool              `json:"allow"`
	Violations []PolicyViolation `json:"violations"`
}
