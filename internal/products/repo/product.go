package repo

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"multitenancypfe/internal/products/dto"
	"multitenancypfe/internal/products/models"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindByID(id, storeID uuid.UUID) (*models.Product, error)
	FindAll(storeID uuid.UUID, filter dto.ProductFilter) ([]models.Product, error)
	Update(product *models.Product) error
	Delete(id, storeID uuid.UUID) error
	SlugExists(slug string, storeID uuid.UUID, excludeID *uuid.UUID) (bool, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id, storeID uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.
		Where("id = ? AND store_id = ? AND deleted_at IS NULL", id, storeID).
		First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &product, err
}

func (r *productRepository) FindAll(storeID uuid.UUID, filter dto.ProductFilter) ([]models.Product, error) {
	var products []models.Product

	query := r.db.Where("store_id = ? AND deleted_at IS NULL", storeID)

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Visibility != nil {
		query = query.Where("visibility = ?", *filter.Visibility)
	}
	if filter.Brand != nil {
		query = query.Where("brand = ?", *filter.Brand)
	}
	if filter.Search != nil && *filter.Search != "" {
		query = query.Where("title ILIKE ?", "%"+strings.TrimSpace(*filter.Search)+"%")
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}

	err := query.
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset((filter.Page - 1) * filter.Limit).
		Find(&products).Error

	return products, err
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id, storeID uuid.UUID) error {
	result := r.db.Model(&models.Product{}).
		Where("id = ? AND store_id = ? AND deleted_at IS NULL", id, storeID).
		Update("deleted_at", gorm.Expr("NOW()"))
	if result.RowsAffected == 0 {
		return errors.New("product not found")
	}
	return result.Error
}

func (r *productRepository) SlugExists(slug string, storeID uuid.UUID, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&models.Product{}).
		Where("slug = ? AND store_id = ? AND deleted_at IS NULL", slug, storeID)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	return count > 0, query.Count(&count).Error
}
