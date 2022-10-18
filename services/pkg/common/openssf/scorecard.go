package openssf

const (
	ScBinaryArtifactsCheck      = "binary_artifacts"
	ScBranchProtectionCheck     = "branch_protection"
	ScCiiBestPracticeCheck      = "cii_best_practices"
	ScCodeReviewCheck           = "code_review"
	ScDangerousWorkflowCheck    = "dangerous_workflow"
	ScDependencyUpdateToolCheck = "dependency_update_tool"
	ScFuzzingCheck              = "fuzzing"
	ScLicenseCheck              = "license"
	ScMaintainedCheck           = "maintained"
	ScPackagingCheck            = "packaging"
	ScPinnedDependenciesCheck   = "pinned_dependencies"
	ScSastCheck                 = "sast"
	ScSecurityPolicyCheck       = "security_policy"
	ScSignedReleasesCheck       = "signed_releases"
	ScVulnerabilitiesCheck      = "vulnerabilities"
	ScTokenPermissionsCheck     = "token_permissions"
)

type ProjectScorecardCheck struct {
	Reason string  `json:"reason"`
	Score  float32 `json:"score"`
}

type ProjectScorecardRepo struct {
	Name   string `json:"name"`
	Commit string `json:"commit"`
}

type ProjectScorecard struct {
	Timestamp uint64                           `json:"timestamp"`
	Score     float32                          `json:"score"`
	Version   string                           `json:"version"`
	Repo      ProjectScorecardRepo             `json:"repo"`
	Checks    map[string]ProjectScorecardCheck `json:"checks"`
}
