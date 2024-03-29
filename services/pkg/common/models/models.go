package models

import (
	config_api "github.com/abhisek/supply-chain-gateway/services/gen"
)

var (
	ArtefactSourceTypeMaven2   = config_api.GatewayUpstreamType_Maven.String()
	ArtefactSourceTypeNpm      = config_api.GatewayUpstreamType_Npm.String()
	ArtefactSourceTypePypi     = config_api.GatewayUpstreamType_PyPI.String()
	ArtefactSourceTypeRubyGems = config_api.GatewayUpstreamType_RubyGems.String()

	ArtefactLicenseTypeSpdx      = "SPDX"
	ArtefactLicenseTypeCycloneDx = "CycloneDX"

	ArtefactVulnerabilitySeverityCritical = "CRITICAL"
	ArtefactVulnerabilitySeverityHigh     = "HIGH"
	ArtefactVulnerabilitySeverityMedium   = "MEDIUM"
	ArtefactVulnerabilitySeverityLow      = "LOW"
	ArtefactVulnerabilitySeverityInfo     = "INFO"

	ArtefactVulnerabilityScoreTypeCVSSv3 = "CVSSv3"
)

type ArtefactRepositoryAuthentication struct {
	Type string `yaml:"type"`
}

type ArtefactUpstreamAuthentication struct {
	Type     string `yaml:"type"`
	Provider string `yaml:"provider"`
}

type ArtefactRepository struct {
	Host           string                           `yaml:"host"`
	Port           int16                            `yaml:"port"`
	Tls            bool                             `yaml:"tls"`
	Sni            string                           `yaml:"sni"`
	Authentication ArtefactRepositoryAuthentication `yaml:"authentication"`
}

type ArtefactRoutingRule struct {
	Prefix string `yaml:"prefix"`
	Host   string `yaml:"host"`
}

type ArtefactUpStream struct {
	Name           string                         `yaml:"name"`
	Type           string                         `yaml:"type"`
	RoutingRule    ArtefactRoutingRule            `yaml:"route"`
	Repository     ArtefactRepository             `yaml:"repository"`
	Authentication ArtefactUpstreamAuthentication `yaml:"authentication"`
}

type ArtefactSource struct {
	Type string `json:"type"`
}

// Align with CVSS v3 but keep room
type ArtefactVulnerabilityScore struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Align with CVE but keep room for enhancement
type ArtefactVulnerabilityId struct {
	Source string `json:"source"`
	Id     string `json:"id"`
}

type ArtefactVulnerability struct {
	Name     string                       `json:"name"`
	Id       ArtefactVulnerabilityId      `json:"id"`
	Severity string                       `json:"severity"`
	Scores   []ArtefactVulnerabilityScore `json:"scores"`
}

// Align with SPDX / CycloneDX
type ArtefactLicense struct {
	Type string `json:"type"` // SPDX | CyloneDX
	Id   string `json:"id"`   // SPDX or CycloneDX ID
	Name string `json:"name"` // Human Readable Name
}

type Artefact struct {
	Source  ArtefactSource `json:"source"`
	Group   string         `json:"group"`
	Name    string         `json:"name"`
	Version string         `json:"version"`
}
