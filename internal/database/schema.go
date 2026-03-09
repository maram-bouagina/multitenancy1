package database

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	productModels "multitenancypfe/internal/products/models"
	storeModel "multitenancypfe/internal/store/models"
)

func CreateTenantSchema(tenantID string) error {
	schema := fmt.Sprintf("tenant_%s", tenantID)

	// 1. Create the schema
	if err := DB.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %q", schema)).Error; err != nil {
		return fmt.Errorf("create schema failed: %w", err)
	}

	// 2. Get the underlying *sql.DB and open a dedicated connection
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB failed: %w", err)
	}

	conn, err := sqlDB.Conn(context.Background())
	if err != nil {
		return fmt.Errorf("get connection failed: %w", err)
	}
	defer conn.Close()

	// 3. Pin search_path to tenant schema on this connection
	if _, err := conn.ExecContext(context.Background(),
		fmt.Sprintf("SET search_path TO %q, public", schema),
	); err != nil {
		return fmt.Errorf("set search_path failed: %w", err)
	}

	// 4. Open a GORM session on that pinned connection
	scopedDB, err := gorm.Open(postgres.New(postgres.Config{Conn: conn}), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open scoped db failed: %w", err)
	}

	// 5. AutoMigrate inside the tenant schema
	if err := scopedDB.AutoMigrate(
		&storeModel.Store{},
		&productModels.Product{},
		&productModels.Category{},
		&productModels.Collection{},
	); err != nil {
		return fmt.Errorf("automigrate failed: %w", err)
	}

	return nil
}
