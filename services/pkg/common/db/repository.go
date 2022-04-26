package db

import (
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db/adapters"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db/models"
)

type VulnerabilityRepository struct {
	adapter adapters.SqlDataAdapter
}

func NewVulnerabilityRepository(adapter adapters.SqlDataAdapter) (*VulnerabilityRepository, error) {
	return &VulnerabilityRepository{adapter: adapter}, nil
}

func (r *VulnerabilityRepository) Upsert(vulnerability models.Vulnerability) error {
	return nil
}

func (r *VulnerabilityRepository) Lookup(ecosystem, group, name string) ([]models.Vulnerability, error) {
	return []models.Vulnerability{}, nil
}
