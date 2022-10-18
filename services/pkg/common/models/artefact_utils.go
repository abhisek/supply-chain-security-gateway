package models

import (
	"fmt"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/openssf"
)

func NewArtefact(src ArtefactSource, name, group, version string) Artefact {
	return Artefact{
		Source:  src,
		Name:    name,
		Group:   group,
		Version: version,
	}
}

func (a Artefact) OpenSsfEcosystem() string {
	if a.Source.Type == ArtefactSourceTypeMaven2 {
		return openssf.VulnerabilityEcosystemMaven
	} else if a.Source.Type == ArtefactSourceTypePypi {
		return openssf.VulnerabilityEcosystemPypi
	} else if a.Source.Type == ArtefactSourceTypeNpm {
		return openssf.VulnerabilityEcosystemNpm
	} else if a.Source.Type == ArtefactSourceTypeRubyGems {
		return openssf.VulnerabilityEcosystemRubyGems
	}

	return ""
}

func (a Artefact) OpenSsfPackageName() string {
	if a.Source.Type == ArtefactSourceTypeMaven2 {
		return fmt.Sprintf("%s:%s", a.Group, a.Name)
	}

	return a.Name
}

func (a Artefact) OpenSsfPackageVersion() string {
	return a.Version
}
