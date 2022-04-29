package models

import (
	common_models "github.com/abhisek/supply-chain-gateway/services/pkg/common/models"
)

const (
	EventTypeVulnerabilityCreated = "event.vulnerability.created"
	EventTypeVulnerabilityUpdated = "event.vulnerability.updated"
)

func NewVulnerabilityDomainEvent(v Vulnerability, t string) common_models.DomainEvent[Vulnerability] {
	return common_models.DomainEvent[Vulnerability]{
		MetaEventWithAttributes: common_models.MetaEventWithAttributes{
			MetaEvent: common_models.MetaEvent{
				Type:    t,
				Version: common_models.EventSchemaVersion,
			},
			MetaAttributes: common_models.MetaAttributes{},
		},
		Data: v,
	}
}
