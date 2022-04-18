package models

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
