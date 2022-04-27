package models

import (
	"fmt"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/openssf"
)

func NewArtefact(src ArtefactSource, name, group, version string) Artefact {
	return Artefact{
		Source:          src,
		Name:            name,
		Group:           group,
		Version:         version,
		Vulnerabilities: []ArtefactVulnerability{},
		Licenses:        []ArtefactLicense{},
	}
}

func (a Artefact) OpenSsfEcosystem() string {
	if a.Source.Type == ArtefactSourceTypeMaven2 {
		return openssf.VulnerabilityEcosystemMaven
	}

	return ""
}

func (a Artefact) OpenSsfPackageName() string {
	if a.Source.Type == ArtefactSourceTypeMaven2 {
		return fmt.Sprintf("%s:%s", a.Group, a.Name)
	}

	return ""
}
