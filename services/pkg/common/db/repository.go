package db

import (
	common_config "github.com/abhisek/supply-chain-gateway/services/pkg/common/config"

	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db/adapters"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db/models"
	"gorm.io/gorm"
)

type VulnerabilityRepository struct {
	config  *common_config.Config
	adapter adapters.SqlDataAdapter
}

func NewVulnerabilityRepository(config *common_config.Config, adapter adapters.SqlDataAdapter) (*VulnerabilityRepository, error) {
	return &VulnerabilityRepository{adapter: adapter}, nil
}

func (r *VulnerabilityRepository) Upsert(vulnerability models.Vulnerability) error {
	db, err := r.adapter.GetDB()
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		var records []models.Vulnerability
		ntx := db.Where(&models.Vulnerability{
			ExternalSource: vulnerability.ExternalSource,
			ExternalId:     vulnerability.ExternalId,
		}).Find(&records)

		if ntx.Error == nil && len(records) > 0 {
			if records[0].DataModifiedAt.Unix() < vulnerability.DataModifiedAt.Unix() {
				vulnerability.ID = records[0].ID
				vulnerability.CreatedAt = records[0].CreatedAt
				ntx = db.Save(&vulnerability)
			}
		} else {
			ntx = db.Create(&vulnerability)
		}

		return ntx.Error
	})

	return err
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
