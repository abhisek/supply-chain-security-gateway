package models

const (
	ArtefactSourceTypeMaven2 = "maven2"
	ArtefactSourceTypeNpm    = "npm"
	ArtefactSourceTypePypi   = "pypi"

	ArtefactLicenseTypeSpdx      = "spdx"
	ArtefactLicenseTypeCycloneDx = "cyclonedx"
)

type ArtefactChannelAuthentication struct{}

type ArtefactRepository struct {
	Host string `yaml:"host"`
	Port int16  `yaml:"port"`
	Tls  bool   `yaml:"tls"`
	Sni  string `yaml:"sni"`
}

type ArtefactRoutingRule struct {
	Prefix string `yaml:"prefix"`
	Host   string `yaml:"host"`
}

type ArtefactUpStream struct {
	Name           string                        `yaml:"name"`
	Type           string                        `yaml:"type"`
	RoutingRule    ArtefactRoutingRule           `yaml:"route"`
	Authentication ArtefactChannelAuthentication `yaml:"authentication"`
	Repository     ArtefactRepository            `yaml:"repository"`
}

type ArtefactSource struct {
	Type string `json:"type"`
}

// Align with CVSS v3 but keep room for enhancement
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
	Name  string                     `json:"name"`
	Id    ArtefactVulnerabilityId    `json:"id"`
	Score ArtefactVulnerabilityScore `json:"score"`
}

// Align with SPDX / CycloneDX
type ArtefactLicense struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Artefact struct {
	Source          ArtefactSource          `json:"source"`
	Group           string                  `json:"group"`
	Name            string                  `json:"name"`
	Version         string                  `json:"version"`
	Vulnerabilities []ArtefactVulnerability `json:"vulnerabilities"`
	Licenses        []ArtefactLicense       `json:"licenses"`
}
