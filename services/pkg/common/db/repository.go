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
	db, err := r.adapter.GetDB()
	if err != nil {
		return err
	}

	tx := db.Create(&vulnerability)
	return tx.Error
}

func (r *VulnerabilityRepository) Lookup(ecosystem, group, name string) ([]models.Vulnerability, error) {
	vulnerabilities := make([]models.Vulnerability, 0)
	db, err := r.adapter.GetDB()
	if err != nil {
		return vulnerabilities, err
	}

	tx := db.Where(&models.Vulnerability{
		Ecosystem: ecosystem,
		Group:     group,
		Name:      name,
	}).Find(&vulnerabilities)

	return vulnerabilities, tx.Error
}
