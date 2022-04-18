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

type PolicyInputVersion struct {
	Major int8 `json:"major"`
	Minor int8 `json:"minor"`
	Patch int8 `json:"patch"`
}

type PolicyInputTarget struct {
	Artefact PolicyEvalTargetArtefact `json:"artefact"`
}

type PolicyInput struct {
	Kind    string             `json:"kind"`
	Version PolicyInputVersion `json:"version"`
	Target  PolicyInputTarget  `json:"target"`
}

type PolicyViolation struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type PolicyResponse struct {
	Allow      bool              `json:"allow"`
	Violations []PolicyViolation `json:"violations"`
}
