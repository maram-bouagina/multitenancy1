package repo

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"multitenancypfe/internal/products/models"
)

type CollectionRepository interface {
	Create(col *models.Collection) error
	FindByID(id, storeID uuid.UUID) (*models.Collection, error)
	FindAll(storeID uuid.UUID) ([]models.Collection, error)
	Update(col *models.Collection) error
	Delete(id, storeID uuid.UUID) error
	AddProduct(collectionID, productID uuid.UUID) error
	RemoveProduct(collectionID, productID uuid.UUID) error
	FindProducts(col *models.Collection, storeID uuid.UUID, page, limit int) ([]models.Product, error)
	SlugExists(slug string, storeID uuid.UUID, excludeID *uuid.UUID) (bool, error)
}

type collectionRepository struct{ db *gorm.DB }

func NewCollectionRepository(db *gorm.DB) CollectionRepository {
	return &collectionRepository{db: db}
}

func (r *collectionRepository) Create(col *models.Collection) error {
	return r.db.Create(col).Error
}

func (r *collectionRepository) FindByID(id, storeID uuid.UUID) (*models.Collection, error) {
	var col models.Collection
	err := r.db.Where("id = ? AND store_id = ?", id, storeID).First(&col).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &col, err
}

func (r *collectionRepository) FindAll(storeID uuid.UUID) ([]models.Collection, error) {
	var cols []models.Collection
	err := r.db.Where("store_id = ?", storeID).Order("name ASC").Find(&cols).Error
	return cols, err
}

func (r *collectionRepository) Update(col *models.Collection) error {
	return r.db.Save(col).Error
}

func (r *collectionRepository) Delete(id, storeID uuid.UUID) error {
	result := r.db.Where("id = ? AND store_id = ?", id, storeID).Delete(&models.Collection{})
	if result.RowsAffected == 0 {
		return errors.New("collection not found")
	}
	return result.Error
}

func (r *collectionRepository) AddProduct(collectionID, productID uuid.UUID) error {
	col := models.Collection{ID: collectionID}
	product := models.Product{ID: productID}
	return r.db.Model(&col).Association("Products").Append(&product)
}

func (r *collectionRepository) RemoveProduct(collectionID, productID uuid.UUID) error {
	col := models.Collection{ID: collectionID}
	product := models.Product{ID: productID}
	return r.db.Model(&col).Association("Products").Delete(&product)
}

// FindProducts handles manual (join) and automatic (rule-based) collections.
// Rule format: "price > 50", "brand = Nike", "stock > 0"
func (r *collectionRepository) FindProducts(col *models.Collection, storeID uuid.UUID, page, limit int) ([]models.Product, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	var products []models.Product

	if col.Type == models.CollectionManual {
		err := r.db.
			Joins("JOIN collection_products cp ON cp.product_id = products.id").
			Where("cp.collection_id = ? AND products.store_id = ? AND products.deleted_at IS NULL", col.ID, storeID).
			Preload("Category").
			Order("products.created_at DESC").
			Limit(limit).Offset((page - 1) * limit).
			Find(&products).Error
		return products, err
	}

	// Automatic: apply rule as WHERE clause
	query := r.db.Where("store_id = ? AND deleted_at IS NULL", storeID)
	if col.Rule != nil && *col.Rule != "" {
		query = applyRule(query, *col.Rule)
	}
	err := query.
		Preload("Category").
		Order("created_at DESC").
		Limit(limit).Offset((page - 1) * limit).
		Find(&products).Error
	return products, err
}

func (r *collectionRepository) SlugExists(slug string, storeID uuid.UUID, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.Collection{}).Where("slug = ? AND store_id = ?", slug, storeID)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	return count > 0, query.Count(&count).Error
}

// applyRule parses simple rule strings like "price > 50", "brand = Nike", "stock > 0"
func applyRule(query *gorm.DB, rule string) *gorm.DB {
	rule = strings.TrimSpace(rule)
	operators := []string{" >= ", " <= ", " > ", " < ", " = "}
	for _, op := range operators {
		if idx := strings.Index(rule, op); idx != -1 {
			field := strings.TrimSpace(rule[:idx])
			value := strings.TrimSpace(rule[idx+len(op):])
			allowed := map[string]bool{
				"price": true, "stock": true, "brand": true,
				"status": true, "visibility": true,
			}
			if allowed[field] {
				query = query.Where(field+" "+strings.TrimSpace(op)+" ?", value)
			}
			break
		}
	}
	return query
}
