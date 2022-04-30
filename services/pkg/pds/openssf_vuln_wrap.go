package pds

import (
	"encoding/json"
	"strings"

	api "github.com/abhisek/supply-chain-gateway/services/gen"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db/models"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/openssf"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/utils"
)

type openssfVulnWrapper struct {
	model                models.Vulnerability
	openssfVulnerability openssf.OsvVulnerability
}

func init() {
	registerVulnerabilitySchemaWrapper(vulnerabilitySchemaWrapperRegistration{
		CanHandle: func(t, v string) bool {
			return (t == models.VulnerabilitySchemaTypeOpenSSF)
		},
		Handle: func(v models.Vulnerability) (VulnerabilitySchemaWrapper, error) {
			w := &openssfVulnWrapper{model: v}
			err := json.Unmarshal(v.Data, &w.openssfVulnerability)

			return w, err
		},
	})
}

func (w *openssfVulnWrapper) CVE() string {
	aliases := utils.SafelyGetValue(w.openssfVulnerability.Aliases)
	for _, alias := range aliases {
		if strings.HasPrefix(alias, "CVE-") {
			return alias
		}
	}

	return ""
}

func (w *openssfVulnWrapper) References() []*api.VulnerabilityReference {
	refs := utils.SafelyGetValue(w.openssfVulnerability.References)
	vRefs := []*api.VulnerabilityReference{}

	for _, i := range refs {
		vRefs = append(vRefs, &api.VulnerabilityReference{
			Type: string(utils.SafelyGetValue(i.Type)),
			Url:  utils.SafelyGetValue(i.Url),
		})
	}

	return vRefs
}

func (w *openssfVulnWrapper) CWEs() []string {
	cwes := w.databaseSpecific()["cwe_ids"]
	if cs, ok := cwes.([]string); ok {
		return cs
	} else {
		return []string{}
	}
}

func (w *openssfVulnWrapper) Affects() []*vulnerabilityAffected {
	av := []*vulnerabilityAffected{}

	osvAffected := utils.SafelyGetValue(w.openssfVulnerability.Affected)
	for _, osvA := range osvAffected {
		av = append(av, &vulnerabilityAffected{
			Versions: utils.SafelyGetValue(osvA.Versions),
		})
	}

	return av
}

func (w *openssfVulnWrapper) FriendlySeverity() string {
	s := w.databaseSpecific()["severity"]
	if cs, ok := s.(string); ok {
		return cs
	} else {
		return "<UNKNOWN>"
	}
}

func (w *openssfVulnWrapper) FriendlySeverityCode() api.VulnerabilitySeverity {
	switch w.FriendlySeverity() {
	case "CRITICAL":
		return api.VulnerabilitySeverity_CRITICAL
	case "HIGH":
		return api.VulnerabilitySeverity_HIGH
	case "MEDIUM":
		return api.VulnerabilitySeverity_HIGH
	case "LOW":
		return api.VulnerabilitySeverity_LOW
	case "INFO":
		return api.VulnerabilitySeverity_INFO
	default:
		return api.VulnerabilitySeverity_UNKNOWN_SEVERITY
	}
}

func (w *openssfVulnWrapper) Severity() []*api.VulnerabilityScore {
	severity := []*api.VulnerabilityScore{}
	s := utils.SafelyGetValue(w.openssfVulnerability.Severity)

	for _, osvSev := range s {
		severity = append(severity, &api.VulnerabilityScore{
			Type:  string(utils.SafelyGetValue(osvSev.Type)),
			Value: utils.SafelyGetValue(osvSev.Score),
		})
	}

	return severity
}

func (w *openssfVulnWrapper) databaseSpecific() map[string]interface{} {
	return utils.SafelyGetValue(w.openssfVulnerability.DatabaseSpecific)
}
