package openssf

const (
	VulnerabilityEcosystemMaven    = "Maven"
	VulnerabilityEcosystemNpm      = "npm"
	VulnerabilityEcosystemPypi     = "PyPI"
	VulnerabilityEcosystemRubyGems = "RubyGems"
)

// https://ossf.github.io/osv-schema/
type OpenSourceVulnerability struct{}

func NewOpenSourceVulnerability(data []byte) (OpenSourceVulnerability, error) {
	return OpenSourceVulnerability{}, nil
}
