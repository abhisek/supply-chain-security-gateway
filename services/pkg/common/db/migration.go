package db

import (
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db/adapters"
	"github.com/abhisek/supply-chain-gateway/services/pkg/common/db/models"
)

func MigrateSqlModels(adapter adapters.SqlDataAdapter) error {
	return adapter.Migrate(&models.Vulnerability{})
}
