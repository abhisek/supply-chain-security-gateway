package models

const (
	ArtefactSourceTypeMaven2 = "maven2"
	ArtefactSourceTypeNpm    = "npm"
	ArtefactSourceTypePypi   = "pypi"
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
	Type string
}

// Align with CVSS v3 but keep room for enhancement
type ArtefactVulnerabilityScore struct {
}

// Align with CVE but keep room for enhancement
type ArtefactVulnerabilityId struct {
}

type ArtefactVulnerability struct {
	Name  string
	Id    ArtefactVulnerabilityId
	Score ArtefactVulnerabilityScore
}

// Align with SPDX / CycloneDX
type ArtefactLicense struct {
}

type Artefact struct {
	Source          ArtefactSource
	Group           string
	Name            string
	Version         string
	Vulnerabilities []ArtefactVulnerability
	Licenses        []ArtefactLicense
}
