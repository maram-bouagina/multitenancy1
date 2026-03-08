package database

import (
	"fmt"
	storeModel "multitenancypfe/internal/store/models"
)

// CreateTenantSchema est appelé uniquement quand un tenant s'inscrit.
// Il crée le schema isolé + les tables dedans.
func CreateTenantSchema(tenantID string) error {
	schema := fmt.Sprintf("tenant_%s", tenantID)

	// 1. Créer le schema
	if err := DB.Exec(fmt.Sprintf(
		"CREATE SCHEMA IF NOT EXISTS %q", schema,
	)).Error; err != nil {
		return fmt.Errorf("create schema failed: %w", err)
	}

	// 2. Pointer vers ce schema
	scopedDB := DB.Exec(fmt.Sprintf("SET search_path TO %q, public", schema))

	// 3. Créer les tables dans ce schema
	if err := scopedDB.AutoMigrate(
		&storeModel.Store{},
	); err != nil {
		return fmt.Errorf("automigrate tenant schema failed: %w", err)
	}

	return nil
}
